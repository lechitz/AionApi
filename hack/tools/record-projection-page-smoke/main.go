package main

import (
	"context"
	"log"
	"os"
)

func main() {
	ctx := context.Background()
	cfg := loadConfig(os.Args[1:])
	if err := run(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}
