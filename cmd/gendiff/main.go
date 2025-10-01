package main

import (
	"code"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "gendiff",
		Usage: "Compares two configuration files and shows a difference.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "output format (default: \"stylish\")",
			},
			&cli.BoolFlag{
				Name:    "help",
				Aliases: []string{"h"},
				Usage:   "show help",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.NArg() != 2 {
				return cli.Exit("Error: Expected 2 file paths", 1)
			}

			filepath1 := cmd.Args().Get(0)
			filepath2 := cmd.Args().Get(1)
			format := cmd.String("format")
			if format == "" {
				format = "stylish"
			}

			result := code.GenDiff(filepath1, filepath2, format)
			fmt.Println(result)

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
