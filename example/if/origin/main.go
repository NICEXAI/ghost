package main

import "fmt"

func main() {
	// @Lazy if:env=="dev" var:msg>hello scope:3
	fmt.Println("current env is dev")
	// @Lazy var:name>world scope:1
	fmt.Println("hello, world")

	fmt.Println("server is running")
}
