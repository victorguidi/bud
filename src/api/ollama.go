package api

import (
	"context"
	"fmt"
	"log"
	"strings"
)

type OllamaAPI struct {
	Url   string
	Model string
	Extensions
}

type Extensions struct {
	Embedder  string
	Streaming bool
}

type OllamaResponse struct {
	Model      string `json:"model"`
	Created_at string `json:"created_at"`
	Response   string `json:"response"`
}

type OllamaEmbeddingResponse struct {
	Embeggind []float64 `json:"embedding"`
}

func New() *OllamaAPI {
	return &OllamaAPI{
		Url:   "http://localhost:11434/api/",
		Model: "llama2",
		Extensions: Extensions{
			Streaming: false,
			Embedder:  "all-minilm",
		},
	}
}

func (o *OllamaAPI) WithUrl(url string) {
	o.Url = url
}

func (o *OllamaAPI) WithModel(modelName string) {
	o.Model = modelName
}

func (o *OllamaAPI) WithEmbedder(embedderName string) {
	o.Embedder = embedderName
}

func (o *OllamaAPI) WithStreaming(stream bool) {
	o.Streaming = stream
}

func (o *OllamaAPI) SendMessageTo(ctx context.Context, msg string) (interface{}, error) {
	apiUrl := o.Url + "generate"
	body := fmt.Sprintf(`{"model": %s, "prompt": %s, "stream": %t}`, o.Model, msg, o.Streaming)
	userData := []byte(body)
	var resp OllamaResponse

	err := HttpCaller(POST, apiUrl, userData, &resp)
	if err != nil {
		log.Panic(err)
	}

	return resp, nil
}

func (o *OllamaAPI) GenerateEmbedding(ctx context.Context, content string) (interface{}, error) {
	apiUrl := o.Url + "embeddings"
	body := fmt.Sprintf(`{"model": "%s", "prompt": "%s"}`, o.Embedder, strings.Trim(content, "\n"))
	log.Println(body)
	userData := []byte(body)
	var resp OllamaEmbeddingResponse

	err := HttpCaller(POST, apiUrl, userData, &resp)
	if err != nil {
		log.Panic(err)
	}

	return resp, nil
}
