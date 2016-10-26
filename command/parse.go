package command

import (
	"io/ioutil"
	"regexp"

	"github.com/nuuls/log"
)

var regex = regexp.MustCompile(`([^{\r\n]+)\s{\s+([^}]+)\s+}`)

func Parse(bs []byte) map[string]*Command {
	matches := regex.FindAllSubmatch(bs, -1)
	cmds := make(map[string]*Command, len(matches))
	for _, match := range matches {
		cmd := &Command{
			Name:  string(match[1]),
			Reply: string(match[2]),
		}
		cmd.Regex = regexp.MustCompile("(?i)^" + cmd.Name)
		cmds[cmd.Name] = cmd
		log.Info("loaded command", cmd.Name)
	}
	return cmds
}

func ParseDir(path string) map[string]*Command {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	cmds := make(map[string]*Command)
	for _, file := range files {
		bs, err := ioutil.ReadFile(path + "/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}
		commands := Parse(bs)
		for key, val := range commands {
			cmds[key] = val
		}
	}
	return cmds
}
