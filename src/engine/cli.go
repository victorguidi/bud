package engine

import "strings"

func (q *Engine) CliArgs(args []string) {
	for i, arg := range args {
		switch strings.ToLower(arg) {
		case "--question", "-q":
			q.Question = args[i+1]
		default:
			continue
		}
	}
}
