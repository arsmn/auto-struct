package autostruct

import (
	"regexp"
	"strconv"
)

var rx = regexp.MustCompile(`(\w+)\((.*?\{.*?\}.*?|[^()]+)\)`)

type command struct {
	list map[string]string
}

func (c command) isCMD(cmd string) bool {
	_, ok := c.list[cmd]
	return ok
}

func (c command) cmd(cmd string) string {
	return c.list[cmd]
}

func (c command) value() string {
	if c.isCMD("value") {
		return c.cmd("value")
	}

	if c.isJSON() {
		return c.json()
	}

	if c.isRepeat() {
		return c.repeat()
	}

	if c.isRune() {
		return c.rune()
	}

	if c.isByte() {
		return c.byte()
	}

	if c.isChannel() {
		return c.channel()
	}

	return ""
}

func (c command) layout() string {
	return c.list["layout"]
}

func (c command) isValueStruct() bool {
	return c.value() == "struct"
}

func (c command) isJSON() bool {
	return c.isCMD("json")
}

func (c command) json() string {
	return c.cmd("json")
}

func (c command) isRepeat() bool {
	return c.isCMD("repeat")
}

func (c command) repeat() string {
	return c.cmd("repeat")
}

func (c command) isRune() bool {
	return c.isCMD("rune")
}

func (c command) rune() string {
	return c.cmd("rune")
}

func (c command) isByte() bool {
	return c.isCMD("byte")
}

func (c command) byte() string {
	return c.cmd("byte")
}

func (c command) isChannel() bool {
	return c.isCMD("chan")
}

func (c command) channel() string {
	return c.cmd("chan")
}

func (c command) len() int {
	i, _ := strconv.Atoi(c.list["len"])
	return i
}

func (c command) cap() int {
	i, _ := strconv.Atoi(c.list["cap"])
	return i
}

func parseTag(tag string) command {
	list := make(map[string]string)

	matches := rx.FindAllStringSubmatch(tag, -1)
	if len(matches) == 0 {
		list["value"] = tag
	} else {
		for _, match := range matches {
			if len(match) == 3 {
				list[match[1]] = match[2]
			}
		}
	}

	return command{list: list}
}
