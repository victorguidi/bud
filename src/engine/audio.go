// This part of the engine handles Text to Audio Inputs
// It uses the bidings from Whisper.cpp. Credits for:
// It also uses the bindings for portaudio. Credits for:

package engine

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	"gitlab.com/bud.git/src/utils"
)

type AudioEngine struct{}

func NewAudioEngine() *AudioEngine {
	return &AudioEngine{}
}

func (a *AudioEngine) Speak(output string) error {
	log.Println(output)
	cmd := exec.Command("spd-say", "-r", "100", "--wait", output)
	if err := cmd.Run(); err != nil {
		log.Printf("Error on running the command: %v", err)
		return err
	}
	return nil
}

func (a *AudioEngine) CallWhisper() (string, error) {
	modelpath := filepath.Join("src", "models", "ggml-base.en.bin")
	samples, err := utils.ReadWav(filepath.Join("samples", "output.wav"))
	if err != nil {
		return "", err
	}

	// Load the model
	model, err := whisper.New(modelpath)
	if err != nil {
		panic(err)
	}
	defer model.Close()

	// Process samples
	context, err := model.NewContext()
	if err != nil {
		panic(err)
	}

	if err := context.Process(samples, nil, nil); err != nil {
		return "", err
	}

	var cmd string
	for {
		segment, err := context.NextSegment()
		if err != nil {
			break
		}
		fmt.Printf("[%6s->%6s] %s\n", segment.Start, segment.End, segment.Text)
		cmd = segment.Text
	}
	return cmd, nil
}
