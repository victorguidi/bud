// This part of the engine handles Text to Audio Inputs
// It uses the bidings from Whisper.cpp. Credits for:
// It also uses the bindings for portaudio. Credits for:

package engine

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	"github.com/gordonklaus/portaudio"
	"gitlab.com/bud.git/src/utils"
)

const (
	RECORDSECONDS = 5
)

type AudioEngine struct {
	AudioChan         chan bool
	AudioResponseChan chan string
}

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

type recorder struct {
	*portaudio.Stream
	buffer []float32
	i      int
}

func (a *AudioEngine) Listen() {
	portaudio.Initialize()
	defer portaudio.Terminate()

	for {
		select {
		case <-a.AudioChan:
			e, err := newRecorder(time.Second * time.Duration(RECORDSECONDS))
			if err != nil {
				log.Println(err)
				continue
			}

			err = e.Start()
			if err != nil {
				log.Println(err)
				continue
			}
			time.Sleep(time.Duration(RECORDSECONDS) * time.Second)
			err = e.Stop()
			if err != nil {
				log.Println(err)
				continue
			}
			output := "cmd.aiff"
			err = utils.SaveToAIFF(output, e.buffer)
			if err != nil {
				log.Println("ERROR CREATING .aiff FILE", err)
				continue
			}
			err = utils.ConvertAiffToWav(output)
			if err != nil {
				log.Println("ERROR CONVERTING .aif to WAV", err)
				continue
			}
			answer, err := a.CallWhisper()
			if err != nil {
				log.Println(err)
				continue
			}

			log.Println("FINAL ANSWER: ", answer)
			a.AudioResponseChan <- answer
			e.Close()
		default:
			// Sleep to avoid busy-wait
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func newRecorder(duration time.Duration) (*recorder, error) {
	h, err := portaudio.DefaultHostApi()
	if err != nil {
		return nil, err
	}
	p := portaudio.LowLatencyParameters(h.DefaultInputDevice, nil)
	p.Input.Channels = 1
	e := &recorder{buffer: make([]float32, int(p.SampleRate*duration.Seconds()))}
	e.Stream, err = portaudio.OpenStream(p, e.processAudio)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *recorder) processAudio(in, out []float32) {
	for i := range in {
		r.buffer[r.i] = in[i]
		r.i++
	}
}

func (a *AudioEngine) STT() error {
	return nil
}

func (a *AudioEngine) CallWhisper() (string, error) {
	defer func() {
		os.Remove(filepath.Join("samples", "cmd.wav"))
	}()
	modelpath := filepath.Join("src", "models", "ggml-base.en.bin")
	samples, err := utils.ReadWav(filepath.Join("samples", "cmd.wav"))
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

	var cmd strings.Builder
	for {
		segment, err := context.NextSegment()
		if err != nil {
			break
		}
		fmt.Printf("[%6s->%6s] %s\n", segment.Start, segment.End, segment.Text)
		if strings.Contains(strings.ToLower(segment.Text), "silence") {
			continue
		}
		cmd.WriteString(segment.Text)
	}
	return cmd.String(), nil
}
