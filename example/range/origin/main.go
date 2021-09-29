package main

import "fmt"

func main() {
	// @Lazy range:count scope:1
	fmt.Println("ha")

	// @Lazy range:count|>num scope:1
	fmt.Println("current counter is num")

	// @Lazy range:data_list|>(name>lazy,age>12) scope:3
	fmt.Println("I am lazy")
	fmt.Println("My age is 12")

}
