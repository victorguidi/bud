package utils

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

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
