package irc

import "strings"

type MsgType int

const (
	MsgPrivmsg MsgType = iota
)

type Message struct {
	Channel     string
	User        string
	DisplayName string
	Me          bool
	Sub         bool
	Mod         bool
	UserType    string
	Type        MsgType
	Text        string
	Tags        map[string]string
}

func Parse(line string) *Message {
	if !strings.HasPrefix(line, "@") {
		return nil
	}
	spl := strings.SplitN(line, " :", 3)
	if len(spl) < 3 {
		return nil
	}
	tags, middle, text := spl[0], spl[1], spl[2]
	m := &Message{
		Text: text,
	}
	m.User, m.Channel = parseMiddle(middle)
	_ = tags
	return m
}

func parseMiddle(middle string) (user, channel string) {
	spl := strings.SplitN(middle, "!", 2)
	if len(spl) < 1 {
		return
	}
	user = spl[0]
	spl = strings.SplitN(middle, "#", 2)
	if len(spl) < 1 {
		return
	}
	channel = spl[1]
	return
}
