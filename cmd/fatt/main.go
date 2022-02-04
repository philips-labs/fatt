package main

import (
	"context"
	"log"

	"github.com/philips-labs/fatt/cmd/fatt/cli"
)

func main() {
	if err := cli.New().ExecuteContext(context.Background()); err != nil {
		log.Fatalf("error during command execution: %v\n", err)
	}
}
