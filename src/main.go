package main

import "gitlab.com/bud.git/src/engine"

func main() {
	engine := engine.New()
	engine.ProcessFiles()
}
