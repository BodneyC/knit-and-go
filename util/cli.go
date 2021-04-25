package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

type IOform int

const (
	ILLEGAL_IOF IOform = iota
	KNIT_IOF
	AST_IOF
)

func toIOform(s string) (IOform, error) {
	if strings.EqualFold(s, "knit") {
		return KNIT_IOF, nil
	} else if strings.EqualFold(s, "ast") {
		return AST_IOF, nil
	} else {
		return ILLEGAL_IOF, fmt.Errorf("Unknown IOform: %s%s", s, StackLine())
	}
}

type CliArgs struct {
	Inform   IOform
	Infile   string
	Jsonfile string
	NoRun    bool
	LogLevel string
	LogTimer bool
}

func ParseCli() (*CliArgs, error) {
	args := &CliArgs{}
	var informStr string
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "inform",
				Aliases:     []string{"inf"},
				Value:       "knit",
				Usage:       "Input file format",
				Destination: &informStr,
			},
			&cli.StringFlag{
				Name:        "ast",
				Value:       "",
				Usage:       "Write parsed .knit to this file as JSON",
				Destination: &args.Jsonfile,
			},
			&cli.BoolFlag{
				Name:        "no-run",
				Aliases:     []string{"norun"},
				Value:       false,
				Usage:       "Prevent the program from running the pattern",
				Destination: &args.NoRun,
			},
			&cli.StringFlag{
				Name:        "log-level",
				Aliases:     []string{"ll"},
				Value:       "",
				Usage:       "Log level (error, info, debug, trace, etc.)",
				Destination: &args.LogLevel,
			},
			&cli.BoolFlag{
				Name:        "timer",
				Value:       false,
				Usage:       "Log time since start of program",
				Destination: &args.LogTimer,
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return fmt.Errorf("No input files given%s", StackLine())
			}

			args.Infile = c.Args().Get(0)

			var err error
			if args.Inform, err = toIOform(informStr); err != nil {
				return fmt.Errorf("%w%s", err, StackLine())
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		return nil, err
	}

	if args.Infile == "" { // Help was ran but didn't exit for some reason
		os.Exit(0)
	}

	return args, nil
}
