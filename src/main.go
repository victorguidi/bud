package main

import "gitlab.com/bud.git/src/engine"

func main() {
	engine := engine.New()
	engine.ProcessFiles()
}

// TODO: Add handlers for PDF and docx
// TODO: Create a simple Crawler For Websites
// TODO: Add SQLite in order to save user data, like sites to craw, other...
// TODO: Simple Frontend Client that is spawned with service (Talk | search docs)
