package main

import (
	"flag"
	"os"

	"github.com/howeyc/fsnotify"
	"github.com/nuuls/log"
	"github.com/nuuls/lulbot/command"
	"github.com/nuuls/lulbot/irc"
)

const (
	commandsPath = "./commands"
)

var noColor = flag.Bool("nocolor", false, "disable log colors")

var (
	cfg  *Config
	chat *irc.Irc
	cmds map[string]*command.Command
)

func init() {
	flag.Parse()
	log.CallerStrLen = 23
	log.AddLogger(&log.Logger{
		Level:  log.LevelDebug,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Color:  !*noColor,
	})

	cfg = LoadConfig("config.json")
}

func main() {
	cmds = command.ParseDir(commandsPath)
	log.Debug(cmds)
	chat = irc.Init(&irc.Config{
		Pass:      cfg.Pass,
		Nick:      cfg.Nick,
		Reconnect: true,
		Server:    cfg.Server,
		Channels:  cfg.Channels,
	})
	go watchCommands()
	for {
		msg := chat.ReadLine()
		log.Info(msg.Channel + "# " + msg.User + ": " + msg.Text)
		go handleMessage(msg)
	}
}

func handleMessage(msg *irc.Message) {
	for _, cmd := range cmds {
		if cmd.Regex.MatchString(msg.Text) {
			if cmd.Script != nil {
				out := cmd.Script.Exec(msg.Channel, msg.User, msg.Text)
				go func(o <-chan string) {
					for line := range o {
						chat.Say(msg.Channel, line)
					}
				}(out)
			} else {
				chat.Say(msg.Channel, cmd.Reply)
			}
		}
	}
}

func watchCommands() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	err = watcher.Watch(commandsPath)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case ev := <-watcher.Event:
			log.Debug(ev)
			commands := command.ParseDir(commandsPath)
			if commands != nil {
				cmds = commands
				log.Info("reloaded all commands")
			}
		case err := <-watcher.Error:
			log.Error(err)
		}
	}
}
