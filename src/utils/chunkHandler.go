package utils

import (
	"errors"
	"strings"
)

type Chunk struct {
	Text string
	Size int
}

func DevideInChunks(text string, chunkSize int) ([]Chunk, error) {
	if chunkSize <= 0 {
		return nil, errors.New("ERROR: CHUNKSIZE MUST BE GREATER THAN ZERO")
	}

	splitted := strings.Split(text, " ")
	var chunks []Chunk

	for i := 0; i < len(splitted); i += chunkSize {
		end := i + chunkSize
		if end > len(splitted) {
			end = len(splitted)
		}

		chunk := Chunk{Size: end - i}
		chunk.splitter(splitted[i:end])
		chunks = append(chunks, chunk)
	}

	return chunks, nil
}

func (c *Chunk) splitter(splitted []string) {
	var text strings.Builder

	for _, word := range splitted {
		if text.Len() > 0 {
			text.WriteString(" ")
		}
		text.WriteString(word)
	}

	c.Text = text.String()
}
