package command

import "regexp"

type Command struct {
	Name  string
	Regex *regexp.Regexp
	Reply string
}
