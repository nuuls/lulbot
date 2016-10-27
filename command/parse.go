package command

import (
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/nuuls/log"
)

const (
	TypeScript = "script"
	TypeText   = "text"
)

var fileRegex = regexp.MustCompile(`([^{\r\n]+)\s{\s+([^}]+)\s+}`)
var bodyRegex = regexp.MustCompile(`(\w+): (?:\((.+)\))? ?(.+)`)

func ParseFile(bs []byte) map[string]*Command {
	matches := fileRegex.FindAllSubmatch(bs, -1)
	cmds := make(map[string]*Command, len(matches))
	for _, match := range matches {
		cmd := ParseBody(string(match[1]), match[2])
		if cmd != nil {
			cmds[cmd.Name] = cmd
			log.Info("loaded command", cmd.Name)
		}
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
		commands := ParseFile(bs)
		for key, val := range commands {
			cmds[key] = val
		}
	}
	return cmds
}

func ParseBody(prefix string, body []byte) *Command {
	cmd := &Command{}
	spl := strings.SplitN(prefix, ":", 2)
	cmd.Name = spl[0]
	if len(spl) > 1 {
		cmd.Regex = regexp.MustCompile(spl[1])
	} else {
		cmd.Regex = regexp.MustCompile("(?i)^" + regexp.QuoteMeta(cmd.Name))
	}
	m := bodyRegex.FindStringSubmatch(string(body))
	if len(m) < 1 {
		log.Warning("no command found in", cmd.Name)
		return nil
	}
	typ := m[1]
	switch typ {
	case TypeScript:
		spl = strings.Split(m[2], " ")
		cmd.Script = &Script{
			Cmd:  spl[0],
			File: m[3],
		}
		if len(spl) > 1 {
			cmd.Script.Args = spl[1:]
		}
	case TypeText:
		cmd.Reply = m[3]
	default:
		log.Error("unsupported type:", typ, "at", cmd.Name)
		return nil
	}
	return cmd
}
