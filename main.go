package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bodneyc/knit-and-go/ast"
	"github.com/bodneyc/knit-and-go/lexer"
	"github.com/bodneyc/knit-and-go/parser"
	"github.com/bodneyc/knit-and-go/util"
	log "github.com/sirupsen/logrus"
)

const (
	SUCCESS_EX int = iota
	GENERIC_EX
	OPTION_EX
	FILESYS_EX
	LEXER_EX
	PARSER_EX
	RUN_EX
)

func configureLogger(logLevelCli string, timings bool) error {
	const (
		LOG_LEVEL_ENV_VAR = "KNIT_LOG_LEVEL"
		DEFAULT_LOG_LEVEL = "info"
	)

	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: !timings,
	})
	log.SetOutput(os.Stdout)

	logLevelStr := DEFAULT_LOG_LEVEL
	if logLevelCli == "" {
		if logLevelEnv, ok := os.LookupEnv(LOG_LEVEL_ENV_VAR); ok {
			logLevelStr = logLevelEnv
		}
	} else {
		logLevelStr = logLevelCli
	}

	logLevel, err := log.ParseLevel(logLevelStr)
	if err != nil {
		return fmt.Errorf("Invalid log level%w", err)
	}

	log.SetLevel(logLevel)

	return nil
}

func main() {
	args, err := util.ParseCli()
	if err != nil {
		log.Trace(err)
		os.Exit(OPTION_EX)
	}

	if err := configureLogger(args.LogLevel, args.LogTimer); err != nil {
		log.Errorf("Failed to set log level%w", err)
		os.Exit(GENERIC_EX)
	}

	log.Info("Starting knit compiler")
	log.Infof("Attempting to open %s", args.Infile)

	file, err := os.Open(args.Infile)
	if err != nil {
		log.Errorf("Couldn't open input file")
		log.Tracef("%s%w", args.Infile, err)
		os.Exit(OPTION_EX)
	}
	defer file.Close()

	var p *parser.Parser

	switch args.Inform {
	case util.KNIT_IOF:
		log.Infof("Parsing input...")

		l := lexer.NewLexer(file)
		p = parser.NewParser(*l)
		err = p.Parse()
		if err != nil {
			log.Errorf("Failed to parse input file%w", err)
			os.Exit(PARSER_EX)
		}

	case util.JSON_IOF:
		log.Info("Reading input JSON")

		fStat, err := file.Stat()
		if err != nil {
			log.Errorf("Couldn't stat input JSON%w", err)
			os.Exit(FILESYS_EX)
		}

		jsonBytes := make([]byte, fStat.Size())
		if _, err := file.Read(jsonBytes); err != nil {
			log.Errorf("Couldn't read input JSON%w", err)
			os.Exit(FILESYS_EX)
		}

		log.Info("Parsing JSON")

		var rootBlockStmt ast.BlockStmt
		if err = json.Unmarshal(jsonBytes, &rootBlockStmt); err != nil {
			log.Errorf("Couldn't parse input JSON%w", err)
			os.Exit(PARSER_EX)
		}

		p = parser.NewParserFromBlockStmt(rootBlockStmt)
	}

	if args.Jsonfile != "" {
		log.Info("Marshalling...")
		rootJson, err := json.MarshalIndent(p.Root, "", "  ")
		if err != nil {
			panic(err)
		}
		log.Info("Marshalling complete")
		log.Info("Writing to file...")
		if err := ioutil.WriteFile(args.Jsonfile, rootJson, 0644); err != nil {
			log.Errorf("Failed to write to root.json : %w", err)
		}
		log.Info("File written")
	}

	if args.NoRun {
		log.Info("No-run option given, exiting...")
		os.Exit(SUCCESS_EX)
	}

	e := ast.NewEngine()
	p.WalkForLocals(e)
	if err := p.WalkForLines(e); err != nil {
		panic(err)
	}

	e.PrintLines()

}
