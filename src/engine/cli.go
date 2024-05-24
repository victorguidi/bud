package engine

import "log"

type BudCLI string

// Constants for common HTTP methods
const (
	ASKBASE BudCLI = "askbase"
)

func (m BudCLI) String() string {
	return string(m)
}

func (q *Engine) CliParser(cmd BudCLI, param string) []byte {
	switch cmd {
	case ASKBASE:
		log.Println(cmd)
		return []byte("WEEEE")
		// q.QuestionChan <- param
	default:
		return nil
	}
}
