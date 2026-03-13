package main

import (
	"fmt"
	"os"
)

func main() {
	cfg := loadConfig(os.Args[1:])
	if err := run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "outbox diagnose failed: %v\n", err)
		os.Exit(1)
	}
}
