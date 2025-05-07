package main

import (
	"fmt"
	"os"

	"qrcodegen/cmd"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "qrcodegen",
		Usage: "Generate QR codes from URLs",
		Commands: []*cli.Command{
			cmd.GenerateCommand(),
			cmd.ServerCommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
