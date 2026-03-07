package main

import "os"

func main() {
	os.Exit(run())
}

func run() int {
	return runWithDeps(newFXApp, os.Getenv, nil)
}
