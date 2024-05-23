package main

import (
	"os"

	"gitlab.com/bud.git/src/engine"
)

func main() {
	args := os.Args
	engine := engine.New()
	engine.CliArgs(args)
	engine.ProcessFiles()
}

// TODO: Add handlers for PDF and docx
// TODO: Create a simple Crawler For Websites
// TODO: Add SQLite in order to save user data, like sites to craw, other...
// TODO: Simple Frontend Client that is spawned with service (Talk | search docs)
// TODO: Add module to convert Audio to Text
