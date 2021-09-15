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
	//tempName := path.Join(strings.ReplaceAll(currentPath, `\`, `/`), "example/origin/test.go")
	//targetName := path.Join(strings.ReplaceAll(currentPath, `\`, `/`), "example/dist/test_1.go")

	options := make(map[string]string)
	options["name"] = "Lazy"
	options["path"] = "dist"
	options["pack_name"] = "testOne"

	//temp, err := lazyTemplate.Parse(tempName, options)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//if err = temp.SaveAsFile(targetName); err != nil {
	//	fmt.Println(err)
	//}

	originFolder := path.Join(strings.ReplaceAll(currentPath, `\`, `/`), "example/origin")
	targetFolder := path.Join(strings.ReplaceAll(currentPath, `\`, `/`), "example/dist")

	if err := lazyTemplate.ParseAll(originFolder, targetFolder, options); err != nil {
		fmt.Println(err)
	}
}
