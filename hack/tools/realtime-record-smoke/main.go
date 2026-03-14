package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	ctx := context.Background()
	cfg := loadConfig(os.Args[1:])

	if err := run(ctx, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "realtime record smoke failed: %v\n", err)
		os.Exit(1)
	}
}
