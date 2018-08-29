package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"log"
)

var iName = flag.String("Type", "", "pass error type name (required)")

func main() {
	flag.Parse()
	if *iName == "" {
		log.Fatal("must be passed type name")
	}

	var tokenset token.Token

	parser.ParseDir()

	fmt.Println("vim-go")
}
