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
	scopeCommand = "scope"
)

func IsLazyCommand(line string) bool {
	line = strings.TrimSpace(line)
	return strings.HasPrefix(line, lazyTag)
}

type Command struct {
	ValCommand []VarCommand
	IfCommand  []IfCommand
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

func ParseLazyCommand(line string) (command Command, err error) {
	var (
		rangeLine     int
		newVarCommand []VarCommand
		newIfCommand  []IfCommand
	)

	if !IsLazyCommand(line) {
		return command, errors.New("invalid lazy command")
	}

	for _, oTag := range strings.Split(line, " ") {
		if oTag == "//" || oTag == " " || oTag == lazyName {
			continue
		}

		//parse scope command
		if strings.HasPrefix(oTag, scopeCommand) {
			vList := strings.Split(oTag, ":")
			if len(vList) != 2 {
				return command, errors.New("invalid range command, error: " + oTag)
			}
			rangeLine, err = strconv.Atoi(vList[1])
			if err != nil {
				return command, err
			}
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
			replaceCommand.Scope = rangeLine
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
				Scope: rangeLine,
			})
			continue
		}
	}

	for _, varOrder := range command.ValCommand {
		varOrder.Scope = rangeLine
		newVarCommand = append(newVarCommand, varOrder)
	}

	for _, ifOrder := range command.IfCommand {
		ifOrder.Scope = rangeLine
		newIfCommand = append(newIfCommand, ifOrder)
	}

	command.ValCommand = newVarCommand
	command.IfCommand = newIfCommand

	return command, nil
}
