package engine

import (
	"log"
	"strings"
)

type BudCLI string

// Constants for common HTTP methods
const (
	ASK     BudCLI = "ask"
	ASKBASE BudCLI = "askbase"
	DIR     BudCLI = "dir"
	QUIT    BudCLI = "quit"
)

func (m BudCLI) String() string {
	return string(m)
}

func (e *Engine) CliParser(cmd string) []byte {
	c := strings.Split(cmd, " ")
	switch c[0] {
	case DIR.String():

		if strings.Contains(c[1], "-s") {
			e.TriggerChan <- Trigger{
				Trigger: DIR.String(),
				Content: DirTrigger{
					Dir: "",
				},
				QuitChan: make(chan bool),
			}
			return []byte("Starting Processing Dir")

		} else if strings.Contains(c[1], "-S") {
			log.Println("TURNING DOWN SERVICE DIRS")
			Workers[DIR.String()].QuitChan <- true
			return []byte("Stopped Processing Dirs")

		} else {
			if len(c) > 2 {
				return []byte("Please provide at least one path to a dir")
			}
			e.TriggerChan <- Trigger{
				Trigger: DIR.String(),
				Content: DirTrigger{
					Dir: strings.Trim(c[2], "\r\n"),
				},
				QuitChan: make(chan bool),
			}
			return []byte("Processing Dir")
		}

	case ASKBASE.String():
		e.TriggerChan <- Trigger{
			Trigger: ASKBASE.String(),
			Content: AskTrigger{
				Question: strings.Join(c[1:], " "),
			},
			QuitChan: make(chan bool),
		}
		return []byte("Processing Question")

	case ASK.String():
		e.TriggerChan <- Trigger{
			Trigger: ASK.String(),
			Content: AskTrigger{
				Question: strings.Join(c[1:], " "),
			},
			QuitChan: make(chan bool),
		}
		return []byte("Processing Question")

	case QUIT.String():
		e.QuitChan <- true
		return []byte("Stopping Bud")

	default:
		return nil
	}
	return nil
}
