package main

import (
	"encoding/json"

	"io/ioutil"
	"os"

	// "github.com/bodneyc/knit-and-go/ast"
	"github.com/bodneyc/knit-and-go/ast"
	"github.com/bodneyc/knit-and-go/lexer"
	"github.com/bodneyc/knit-and-go/parser"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	// log.SetLevel(log.DebugLevel)
}

func main() {
	log.SetOutput(os.Stderr)

	file, e := os.Open("./test-patterns/rsc/comfy-raglan.knit")
	if e != nil {
		panic(e)
	}
	defer file.Close()

	l := lexer.NewLexer(file)
	p := parser.NewParser(*l)
	e = p.Parse()
	if e != nil {
		log.Error(e)
	}

	log.Info("Marshalling...")
	rootJson, e := json.MarshalIndent(p.Root, "", " ")
	if e != nil {
		panic(e)
	}
	log.Info("Marshalling complete")

	log.Info("Writing to file...")
	if e := ioutil.WriteFile("out/root.json", rootJson, 0644); e != nil {
		log.Errorf("Failed to write to root.json : %w", e)
	}
	log.Info("File written")

	//----

	var readBlock ast.BlockStmt
	if e := json.Unmarshal(rootJson, &readBlock); e != nil {
		panic(e)
	}

	log.Info("Marshalling...")
	readJson, e := json.MarshalIndent(readBlock, "", " ")
	if e != nil {
		panic(e)
	}
	log.Info("Marshalling complete")

	log.Info("Writing to file...")
	if e := ioutil.WriteFile("out/read.json", readJson, 0644); e != nil {
		log.Errorf("Failed to write to root.json : %w", e)
	}
	log.Info("File written")

}
