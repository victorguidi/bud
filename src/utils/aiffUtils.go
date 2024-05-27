package utils

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-audio/aiff"
	"github.com/go-audio/audio"
)

const (
	SAMPLES       = "samples"
	SAMPLERATE    = 44100
	NUMCHANNELS   = 1
	BITSPERSAMPLE = 16
)

func SaveToAIFF(path string, buffer []float32) error {
	log.Println("HERE")
	file, err := os.Create(filepath.Join("samples", path))
	if err != nil {
		return err
	}
	defer file.Close()

	// Convert buffer to int format as AIFF format expects integer samples
	intBuffer := make([]int, len(buffer))
	for i, sample := range buffer {
		intBuffer[i] = int(sample * 32767) // Convert float32 sample to int sample
	}

	// AIFF Encoder setup
	enc := aiff.NewEncoder(file, SAMPLERATE, BITSPERSAMPLE, NUMCHANNELS)
	defer enc.Close()

	buf := &audio.IntBuffer{
		Format: &audio.Format{
			SampleRate:  SAMPLERATE,
			NumChannels: NUMCHANNELS,
		},
		Data: intBuffer,
	}
	err = enc.Write(buf)
	if err != nil {
		return err
	}

	return nil
}

func ConvertAiffToWav(file string) error {
	samples := "samples"
	defer func() {
		os.Remove(filepath.Join(samples, file))
	}()
	newFile := strings.ReplaceAll(file, ".aiff", ".wav")

	// Command that I am running ffmpeg -i input.aiff -ar 16000 -ac 1 -c:a pcm_s16le output.wav
	cmd := exec.Command("ffmpeg",
		"-i",
		filepath.Join(samples, file),
		"-ar",
		"16000",
		"-ac",
		"1",
		"-c:a",
		"pcm_s16le",
		filepath.Join(samples, newFile))

	if err := cmd.Run(); err != nil {
		log.Printf("Error on convert pdf to txt: %v", err)
		return err
	}
	return nil
}
