package main

import (
	"fmt"
)

func main() {

	// @Lazy range:data_list|>(name>ghost,age>12) scope:3
	fmt.Println("I am ghost")
	fmt.Println("My age is 12")

	// @Lazy range:data_list>ghost scope:1
	fmt.Println(`hello, ghost`)
}
