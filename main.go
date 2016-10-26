package main

import (
	"github.com/nuuls/log"
	"github.com/nuuls/lulbot/command"
	"github.com/nuuls/lulbot/irc"
)

var (
	cfg  *Config
	chat *irc.Irc
	cmds map[string]*command.Command
)

func init() {
	log.AddLogger(log.DefaultLogger)
	log.CallerStrLen = 23
	cfg = LoadConfig("config.json")
}

func main() {
	cmds = command.ParseDir("./commands")
	log.Debug(cmds)
	chat = irc.Init(&irc.Config{
		Pass:      cfg.Pass,
		Nick:      cfg.Nick,
		Reconnect: true,
		Server:    cfg.Server,
		Channels:  cfg.Channels,
	})
	for {
		msg := chat.ReadLine()
		log.Info(msg.Channel + "# " + msg.User + ": " + msg.Text)
		go handleMessage(msg)
	}
}

func handleMessage(msg *irc.Message) {
	for _, cmd := range cmds {
		if cmd.Regex.MatchString(msg.Text) {
			chat.Say(msg.Channel, cmd.Reply)
		}
	}
}
