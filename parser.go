package lazyTemplate

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
)

func isLazyCommand(line string) bool {
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
	Range    int    // affected range
}

type IfCommand struct {
	Expr  string // judgment condition
	Range int    // affected range
}

func parseLazyCommand(line string) (command Command, err error) {
	var (
		rangeLine int
		newVarCommand []VarCommand
		newIfCommand []IfCommand
	)

	if !isLazyCommand(line) {
		return command, errors.New("invalid lazy command")
	}

	for _, oTag := range strings.Split(line, " ") {
		if oTag == "//" || oTag == " " || oTag == lazyName {
			continue
		}

		//parse range command
		if strings.HasPrefix(oTag, rangeCommand) {
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
			replaceCommand.Range = rangeLine
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
				Expr: vList[1],
				Range: rangeLine,
			})
			continue
		}
	}

	for _, varOrder := range command.ValCommand {
		varOrder.Range = rangeLine
		newVarCommand = append(newVarCommand, varOrder)
	}

	for _, ifOrder := range command.IfCommand {
		ifOrder.Range = rangeLine
		newIfCommand = append(newIfCommand, ifOrder)
	}

	command.ValCommand = newVarCommand
	command.IfCommand = newIfCommand

	return command, nil
}
