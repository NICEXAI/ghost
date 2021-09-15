package main

// lazy replace:path>origin replace:pack_name>test range:1
import "github.com/NICEXAI/lazy-template-engine/example/origin/test"

func main() {
	// lazy replace:pack_name>test range:1
	test.Test()
}
