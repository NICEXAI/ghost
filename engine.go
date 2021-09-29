package ghost

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NICEXAI/ghost/parser"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/NICEXAI/ghost/util"

	"github.com/fatih/color"
)

const rangeTag = "// @Tag %v_%v"

type Template struct {
	builder strings.Builder
}

func (t *Template) SaveAsFile(name string) error {
	return util.CreateIfNotExist(name, t.builder.String())
}

// Parse parse template files
func Parse(tempName string, options map[string]interface{}) (temp Template, err error) {
	if !util.IsFileExist(tempName) {
		return temp, errors.New("template file not exist")
	}

	f, err := os.Open(tempName)
	defer f.Close()

	if err != nil {
		return temp, err
	}

	buf := bufio.NewReader(f)
	command := parser.Command{}

	for {
		var (
			line             []byte
			newCommand       parser.Command
			nextIfCommand    []parser.IfCommand
			nextVarCommand   []parser.VarCommand
			nextRangeCommand []parser.RangeCommand
			ignoreStatus     bool
		)

		line, _, err = buf.ReadLine()
		if err != nil {
			break
		}

		lineCon := string(line)

		//parse template file
		if parser.IsLazyCommand(lineCon) {
			newCommand, err = parser.ParseLazyCommand(lineCon)
			if err != nil {
				break
			}

			command.IfCommand = append(command.IfCommand, newCommand.IfCommand...)
			command.ValCommand = append(command.ValCommand, newCommand.ValCommand...)

			// preprocess range command
			for _, rangeCommand := range newCommand.RangeCommand {
				rangeCommand, err = parser.ParseAndExecuteRangeExpr(rangeCommand, options)
				if err != nil {
					return temp, err
				}
				command.RangeCommand = append(command.RangeCommand, rangeCommand)
			}

			bData, _ := json.Marshal(command)
			fmt.Println(string(bData))
			continue
		}

		//execute range command
		if len(command.RangeCommand) > 0 {
			for _, order := range command.RangeCommand {

				if order.TagId == 0 {
					order.TagId = rand.Int()

					for i := 0; i < order.Loop; i++ {
						for s := 0; s < order.Scope; s++ {
							temp.builder.WriteString(fmt.Sprintf(rangeTag, order.TagId, s+1) + "\n")
						}
					}
				}

				// replace target content
				order.Counter += 1
				rangeTagCon := fmt.Sprintf(rangeTag, order.TagId, order.Counter)

				for i := 0; i < order.Loop; i++ {
					if len(order.Action) > 0 {
						newLine := lineCon

						for _, action := range order.Action[0][i] {
							newLine = strings.ReplaceAll(newLine, action.Target, action.Value)
						}

						newTempCon := strings.Replace(temp.builder.String(), rangeTagCon, newLine, 1)
						temp.builder.Reset()
						temp.builder.WriteString(newTempCon)
					} else {
						newTempCon := strings.Replace(temp.builder.String(), rangeTagCon, lineCon, 1)
						temp.builder.Reset()
						temp.builder.WriteString(newTempCon)
					}

					fmt.Println(temp.builder.String())
				}

				if len(order.Action) > 0 {
					order.Action = order.Action[1:]
				}

				if order.Counter < order.Scope {
					nextRangeCommand = append(nextRangeCommand, order)
				}
			}

			command.RangeCommand = nextRangeCommand
			continue
		}

		//execute if command
		if len(command.IfCommand) > 0 {
			for _, order := range command.IfCommand {
				var res interface{}

				if order.Scope > 0 {
					if !ignoreStatus {
						res, err = parser.ParseAndExecuteExpr(order.Expr, options)
						if err != nil {
							return temp, err
						}

						bRes, ok := res.(bool)
						if !ok {
							return temp, errors.New("expr must be a bool")
						}

						if !bRes {
							ignoreStatus = true
						}
					}

					order.Scope--
				}
				if order.Scope > 0 {
					nextIfCommand = append(nextIfCommand, order)
				}
			}

			command.IfCommand = nextIfCommand
		}

		//execute var command
		if !ignoreStatus && len(command.ValCommand) > 0 {
			for _, order := range command.ValCommand {
				if order.Scope > 0 {
					var (
						originVal interface{}
						lData     []byte
						lastVal   string
					)

					originVal = options[order.Variable]
					switch lv := originVal.(type) {
					case string:
						lastVal = lv
					case int:
						lastVal = strconv.Itoa(lv)
					case map[string]interface{}:
						lData, err = json.Marshal(lv)
						if err != nil {
							return Template{}, err
						}
						lastVal = string(lData)
					}
					lineCon = strings.ReplaceAll(lineCon, order.Target, lastVal)
					order.Scope--
				}
				if order.Scope > 0 {
					nextVarCommand = append(nextVarCommand, order)
				}
			}

			command.ValCommand = nextVarCommand
		}

		if !ignoreStatus {
			temp.builder.WriteString(lineCon + "\n")
		}
	}

	return temp, nil
}

// ParseAll parse all template files in the folder
func ParseAll(originFolder, targetFolder string, options map[string]interface{}) error {
	if !util.IsFolderExist(originFolder) {
		return errors.New("origin folder is not exist")
	}

	if !util.IsFolderExist(targetFolder) {
		return errors.New("target folder is not exist")
	}

	files, err := util.GetFileListFromFolder(originFolder)
	if err != nil {
		fmt.Println(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(len(files))

	for _, file := range files {
		go func(fileInfo util.FileInfo) {
			defer wg.Done()

			temp, err := Parse(fileInfo.Path, options)
			if err != nil {
				color.Red("file %s parse failed, error: %s", fileInfo.Name, err.Error())
				return
			}

			if err := temp.SaveAsFile(strings.ReplaceAll(fileInfo.Path, originFolder, targetFolder)); err != nil {
				color.Red("file %s save failed, error: %s", fileInfo.Name, err.Error())
			}
		}(file)
	}

	wg.Wait()

	return nil
}
