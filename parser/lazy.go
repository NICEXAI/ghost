package parser

import (
	"errors"
	"strconv"
	"strings"
)

const (
	lazyTag  = "// @Lazy"
	lazyName = "@Lazy"

	varCommand   = "var"
	ifCommand    = "if"
	rangeCommand = "range"
	scopeCommand = "scope"
)

func IsLazyCommand(line string) bool {
	line = strings.TrimSpace(line)
	return strings.HasPrefix(line, lazyTag)
}

type Command struct {
	ValCommand   []VarCommand
	IfCommand    []IfCommand
	RangeCommand []RangeCommand
}

type VarCommand struct {
	Variable string // variable
	Target   string // replace the content
	Scope    int    // affected scope
}

type IfCommand struct {
	Expr  string // judgment condition
	Scope int    // affected scope
}

type RangeCommand struct {
	Expr    string // judgment condition
	Scope   int    // affected scope
	Loop    int
	TagId   int // slot tag id
	Counter int
	Action  [][][]RangeAction
}

type RangeAction struct {
	Target string // replace the content
	Value  string // current slot value
}

func ParseLazyCommand(line string) (command Command, err error) {
	var (
		scopeLine int
	)

	if !IsLazyCommand(line) {
		return command, errors.New("invalid lazy command")
	}

	oTagList := strings.Split(line, " ")

	for _, oTag := range oTagList {
		//parse scope command
		if strings.HasPrefix(oTag, scopeCommand) {
			vList := strings.Split(oTag, ":")
			if len(vList) != 2 {
				return command, errors.New("invalid range command, error: " + oTag)
			}
			scopeLine, err = strconv.Atoi(vList[1])
			if err != nil {
				return command, err
			}
			break
		}
	}

	for _, oTag := range oTagList {
		if oTag == "//" || oTag == " " || oTag == lazyName {
			continue
		}

		//parse var command
		if strings.HasPrefix(oTag, varCommand) {
			vList := strings.Split(oTag, ">")
			if len(vList) != 2 {
				return command, errors.New("invalid var command, error: " + oTag)
			}
			replaceCommand := VarCommand{}
			replaceCommand.Variable = vList[0][len(varCommand)+1:]
			replaceCommand.Target = vList[1]
			replaceCommand.Scope = scopeLine
			command.ValCommand = append(command.ValCommand, replaceCommand)
			continue
		}

		//parse if command
		if strings.HasPrefix(oTag, ifCommand) {
			vList := strings.Split(oTag, ":")
			if len(vList) != 2 {
				return command, errors.New("invalid if command, error: " + oTag)
			}

			command.IfCommand = append(command.IfCommand, IfCommand{
				Expr:  vList[1],
				Scope: scopeLine,
			})
			continue
		}

		//parse range command
		if strings.HasPrefix(oTag, rangeCommand) {
			vList := strings.Split(oTag, ":")
			if len(vList) != 2 {
				return command, errors.New("invalid range command, error: " + oTag)
			}

			command.RangeCommand = append(command.RangeCommand, RangeCommand{
				Expr:  vList[1],
				Scope: scopeLine,
			})
			continue
		}
	}

	return command, nil
}
