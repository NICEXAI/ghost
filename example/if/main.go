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
	options["env"] = "dev"


	originFolder := path.Join(strings.ReplaceAll(currentPath, `\`, `/`), "example/if/origin")
	targetFolder := path.Join(strings.ReplaceAll(currentPath, `\`, `/`), "example/if/dist")

	if err := lazyTemplate.ParseAll(originFolder, targetFolder, options); err != nil {
		fmt.Println(err)
	}
}
