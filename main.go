package main

// see https://blog.golang.org/generate
//go:generate go run mime/generate.go
//go:generate go fmt mime/mime.go

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	(&cli.App{}).Run(os.Args)
}
