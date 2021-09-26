package lazyTemplate

import (
	"errors"
	"strconv"
	"strings"
)

const (
	lazyTag = "// @Lazy"

	lazyName = "@Lazy"
	replace  = "replace"
	ranges   = "range"
)

func isLazyCommand(line string) bool {
	line = strings.TrimSpace(line)
	return strings.HasPrefix(line, lazyTag)
}

type Command struct {
	Replace []ReplaceCommand
	Range   int // affected range
}

type ReplaceCommand struct {
	Variable string // variable
	Target   string // replace the content
}

func parseLazyCommand(line string) (command Command, err error) {
	if !isLazyCommand(line) {
		return command, errors.New("invalid lazy command")
	}

	for _, oTag := range strings.Split(line, " ") {
		if oTag == "//" || oTag == " " || oTag == lazyName {
			continue
		}

		//parse replace command
		if strings.HasPrefix(oTag, replace) {
			vList := strings.Split(oTag, ">")
			if len(vList) != 2 {
				return command, errors.New("invalid replace command, error: " + oTag)
			}
			replaceCommand := ReplaceCommand{}
			replaceCommand.Variable = vList[0][8:]
			replaceCommand.Target = vList[1]
			command.Replace = append(command.Replace, replaceCommand)
			continue
		}

		//parse range command
		if strings.HasPrefix(oTag, ranges) {
			vList := strings.Split(oTag, ":")
			if len(vList) != 2 {
				return command, errors.New("invalid range command, error: " + oTag)
			}
			rangeNum := 0
			rangeNum, err = strconv.Atoi(vList[1])
			if err != nil {
				return command, err
			}
			command.Range = rangeNum
			continue
		}
	}

	return command, nil
}
