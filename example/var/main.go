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
	options["name"] = "Lazy"
	options["path"] = "dist"
	options["pack_name"] = "testOne"
	options["env"] = "dev"

	originFolder := path.Join(strings.ReplaceAll(currentPath, `\`, `/`), "example/var/origin")
	targetFolder := path.Join(strings.ReplaceAll(currentPath, `\`, `/`), "example/var/dist")

	if err := ghost.ParseAll(originFolder, targetFolder, options); err != nil {
		fmt.Println(err)
	}
}
