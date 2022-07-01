package main

import (
	"log"
	"os"
)

// main call with path, command, arguments.
func main() {
	if len(os.Args) < 3 {
		log.Fatal("path and command must be specified")
	}
	env, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatalf("failed reading dir: %v", err)
	}

	code := RunCmd(os.Args[2:], env)
	os.Exit(code)
}
