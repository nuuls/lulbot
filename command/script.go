package command

import (
	"bufio"
	"os/exec"

	"github.com/nuuls/log"
)

type Script struct {
	Cmd  string
	Args []string
	File string
}

func (s *Script) Exec(channel, user, msg string) <-chan string {
	args := append(s.Args, s.File, channel, user, msg)
	c := exec.Command(s.Cmd, args...)
	stdout, err := c.StdoutPipe()
	if err != nil {
		log.Error(err)
		return nil
	}
	stderr, err := c.StderrPipe()
	if err != nil {
		log.Error(err)
		return nil
	}
	err = c.Start()
	if err != nil {
		log.Error(err)
	}
	out := make(chan string)
	go func() {
		defer close(out)
		reader := bufio.NewReader(stdout)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			out <- line
		}
	}()
	go func() {
		reader := bufio.NewReader(stderr)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			log.Info(s.File, ":", line)
		}
	}()
	return out
}
