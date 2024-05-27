package main

import (
	"gitlab.com/bud.git/src/engine"
)

func main() {
	bud := engine.New()

	go bud.Run()
	go bud.StartServer()

	api := engine.NewBudAPI(bud)
	api.RegisterHandlers()
	api.Start("9875")
}

// TODO: Implement toggle for calling Bud with Audio
// TODO: Add help command
// TODO: Create a simple Crawler For Websites
// TODO: Simple Frontend Client that is spawned with service (Talk | search docs)
// TODO: Implement Integration with Audio Input to capture Microphone
// TODO: Add module to convert Audio to Text
