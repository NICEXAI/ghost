package main

import (
	"fmt"
	"github.com/NICEXAI/ghost"
	"os"
	"path"
	"strings"
)

func main() {
	currentPath, _ := os.Getwd()

	options := make(map[string]interface{})

	options["env"] = "dev"
	options["name"] = "ghost"

	originFolder := path.Join(strings.ReplaceAll(currentPath, `\`, `/`), "example/if/origin")
	targetFolder := path.Join(strings.ReplaceAll(currentPath, `\`, `/`), "example/if/dist")

	if err := ghost.ParseAll(originFolder, targetFolder, options); err != nil {
		fmt.Println(err)
	}
}
