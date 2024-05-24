package main

import (
	"os"

	"gitlab.com/bud.git/src/api"
	"gitlab.com/bud.git/src/engine"
)

func main() {
	args := os.Args
	engine := engine.New()
	engine.CliArgs(args)
	// engine.ProcessFiles()

	api := api.NewBudAPI()
	api.RegisterHandlers()
	api.Start("5000")
}

// TODO: Create a simple Crawler For Websites
// TODO: Add SQLite in order to save user data, like sites to craw, other...
// TODO: Simple Frontend Client that is spawned with service (Talk | search docs)
// TODO: Implement Integration with Audio Input to capture Microphone
// TODO: Add module to convert Audio to Text
