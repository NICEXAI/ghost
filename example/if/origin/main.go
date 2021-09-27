package main

import "fmt"

func main() {
	// @Lazy if:env=="dev" range:3
	fmt.Println("current env is dev")
	// @Lazy var:name>world range:1
	fmt.Println("hello, world")

	fmt.Println("server is running")
}
