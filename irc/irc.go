package irc

import (
	"bufio"
	"crypto/tls"
	"net"
	"net/textproto"
	"strings"
	"time"

	"github.com/nuuls/log"
)

var cfg *Config

type Config struct {
	Pass      string
	Nick      string
	Reconnect bool
	RateLimit bool
	Server    string
	Channels  []string
}

type Irc struct {
	conn     net.Conn
	channels []string
	rate     int
	messages chan *Message
}

func Init(config *Config) *Irc {
	cfg = config
	i := &Irc{
		channels: cfg.Channels,
		messages: make(chan *Message, 10),
	}
	if !strings.HasPrefix(config.Pass, "oauth:") {
		config.Pass = "oauth:" + config.Pass
	}
	i.connect()
	return i
}

func (i *Irc) connect() {
	var err error
	if cfg.Server == "" {
		i.conn, err = tls.Dial("tcp", "irc.chat.twitch.tv:443", nil)
	} else {
		i.conn, err = net.Dial("tcp", cfg.Server)
	}
	if err != nil {
		log.Critical(err)
		time.Sleep(time.Second * 2)
		i.connect()
		return
	}
	go i.read()
	log.Info("connected to chat server")
	err = i.Send("PASS " + cfg.Pass)
	if err != nil {
		log.Error(err)
		i.connect()
		return
	}
	i.Send("NICK " + cfg.Nick)
	i.Send("CAP REQ twitch.tv/tags")
	i.Send("CAP REQ twitch.tv/commands")
	for _, channel := range cfg.Channels {
		i.Join(channel)
	}
}

func (i *Irc) Send(msg string) error {
	_, err := i.conn.Write([]byte(msg + "\r\n"))
	return err
}

func (i *Irc) Join(channel string) {
	if strings.HasPrefix(channel, "#") {
		channel = channel[1:]
	}
	err := i.Send("JOIN #" + channel)
	if err != nil {
		log.Error(err)
		time.AfterFunc(time.Second*5, func() {
			i.Join(channel)
		})
		return
	}
	log.Info("joined channel", channel)
}

func (i *Irc) Say(channel, message string) {
	if i.rate > 17 {
		log.Warning("rate limited")
		return
	}
	err := i.Send("PRIVMSG #" + channel + " :" + message + " ")
	if err != nil {
		log.Error("error sending message:", message, err)
		i.connect()
		return
	}
	log.Info("sent: ", channel, ":", message)
	i.rate++
	time.AfterFunc(time.Second*32, func() {
		i.rate--
	})
}

func (i *Irc) read() {
	tp := textproto.NewReader(bufio.NewReader(i.conn))
	for {
		line, err := tp.ReadLine()
		if err != nil {
			log.Error(err)
			i.connect()
			return
		}
		if strings.HasPrefix(line, "PING") {
			i.Send(strings.Replace(line, "PING", "PONG", 1))
			continue
		}
		msg := Parse(line)
		if msg != nil {
			select {
			case i.messages <- msg:
			default:
				log.Warning("no listener", msg)
			}
		} else {
			log.Debug(line)
		}
	}
}

func (i *Irc) ReadLine() *Message {
	return <-i.messages
}
