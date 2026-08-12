package main

import (
	"os"

	"ddd-bootstrap/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
