package autostruct

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var rx = regexp.MustCompile(`(\w+)(?:\(([^)]*)\))?`)

type command struct {
	txt, val string
	len, cap int
	layout   string
}

func (c command) isJSON() bool {
	return strings.EqualFold(c.txt, "json")
}

func (c command) isRepeat() bool {
	return strings.EqualFold(c.txt, "repeat")
}

func (c command) isRune() bool {
	return strings.EqualFold(c.txt, "rune")
}

func (c command) isValStruct() bool {
	return strings.EqualFold(c.val, "struct")
}

func parseTag(tag string) (command, error) {
	cmd := command{
		val: tag,
	}

	for _, match := range rx.FindAllStringSubmatch(tag, -1) {
		if len(match) == 3 {
			switch match[1] {
			case "value", "json", "repeat", "rune", "byte":
				if cmd.txt != "" {
					return command{}, fmt.Errorf("more than one main command is detected")
				}

				cmd.txt, cmd.val = match[1], match[2]
			case "len":
				l, _ := strconv.Atoi(match[2])
				cmd.len = l
			case "cap":
				l, _ := strconv.Atoi(match[2])
				cmd.cap = l
			case "layout":
				cmd.layout = match[2]
			}
		}
	}

	if cmd.cap < cmd.len {
		cmd.cap = cmd.len
	}

	return cmd, nil
}
