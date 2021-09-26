package main

import (
	"fmt"
	lazyTemplate "github.com/NICEXAI/lazy-template-engine"
	"os"
	"path"
	"strings"
)

func main() {
	currentPath, _ := os.Getwd()

	options := make(map[string]string)
	options["name"] = "Lazy"
	options["path"] = "dist"
	options["pack_name"] = "testOne"
	options["env"] = "dev"


	originFolder := path.Join(strings.ReplaceAll(currentPath, `\`, `/`), "example/var/origin")
	targetFolder := path.Join(strings.ReplaceAll(currentPath, `\`, `/`), "example/var/dist")

	if err := lazyTemplate.ParseAll(originFolder, targetFolder, options); err != nil {
		fmt.Println(err)
	}
}
