package main

// @Lazy replace:path>origin replace:pack_name>test range:1
import "github.com/NICEXAI/lazy-template-engine/example/origin/test"

func main() {
	// @Lazy replace:pack_name>test range:1
	test.Test()
}
