// This part of the engine handles Text to Audio Inputs
package engine

import (
	"log"
	"os/exec"
)

type AudioEngine struct{}

func NewAudioEngine() *AudioEngine {
	return &AudioEngine{}
}

func (a *AudioEngine) Speak(output string) error {
	// spd-say "Arch Linux is the best"
	log.Println(output)
	cmd := exec.Command("spd-say", "--wait", output)
	if err := cmd.Run(); err != nil {
		log.Printf("Error on convert pdf to txt: %v", err)
		return err
	}
	return nil
}
