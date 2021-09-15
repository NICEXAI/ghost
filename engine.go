package lazyTemplate

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/NICEXAI/lazy-template-engine/util"

	"github.com/fatih/color"
)

type Template struct {
	builder strings.Builder
}

func (t *Template) SaveAsFile(name string) error {
	return util.CreateIfNotExist(name, t.builder.String())
}

// Parse parse template files
func Parse(tempName string, options map[string]string) (temp Template, err error) {
	if !util.IsFileExist(tempName) {
		return temp, errors.New("template file not exist")
	}

	f, err := os.Open(tempName)
	defer f.Close()

	if err != nil {
		return temp, err
	}

	buf := bufio.NewReader(f)
	command := Command{}

	for {
		var line []byte

		line, _, err = buf.ReadLine()
		if err != nil {
			break
		}

		lineCon := string(line)

		//parse template file
		if isLazyCommand(lineCon) {
			command, err = parseLazyCommand(lineCon)
			if err != nil {
				break
			}
			continue
		}

		if command.Range > 0 && len(command.Replace) > 0 {
			for _, order := range command.Replace {
				lineCon = strings.ReplaceAll(lineCon, order.Target, options[order.Variable])
			}
			command.Range--
		}

		temp.builder.WriteString(lineCon + "\n")
	}

	return temp, nil
}

// ParseAll parse all template files in the folder
func ParseAll(originFolder, targetFolder string, options map[string]string) error {
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
