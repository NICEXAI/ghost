package main

import (
	"fmt"
	"github.com/NICEXAI/ghost/parser"
)

func main() {
	attrs := make(map[string]interface{})
	attrs["env"] = "test"
	result, err := parser.ParseAndExecuteExpr(`env=="dev"`, attrs)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
}
