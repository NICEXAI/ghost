package main

// @Lazy var:path>origin scope:1
import (
	"github.com/NICEXAI/ghost/example/var/origin/test"
)

func main() {
	// @Lazy var:pack_name>test scope:1
	test.Test()
}
