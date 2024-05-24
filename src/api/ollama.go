package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// NOTE: all-minilm = 384 dimensions
// NOTE: mxbai-embed-large = 1024 demensions
// NOTE: nomic-embed-text = 768

const (
	DEFAULTPROMPT    = "You are a helpfull assistant. Please answer the question provided in the PROMPT."
	DEFAULTRAGPROMPT = "You are a helpfull assistant that provides answer based on the knowledge given to you.If There is no context, answer: I don't know, maybe I need more context."
)

type OllamaAPI struct {
	Url   string
	Model string
	Extensions
}

type Extensions struct {
	Embedder  string
	Prompt    string
	Streaming bool
	Tokens    int
}

type PromptFormater func(map[string]any) string

type OllamaResponse struct {
	Model      string `json:"model"`
	Created_at string `json:"created_at"`
	Response   string `json:"response"`
}

type OllamaEmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

func NewOllamaAPI() *OllamaAPI {
	return &OllamaAPI{
		Url:   "http://localhost:11434/api/",
		Model: "llama2",
		Extensions: Extensions{
			Streaming: false,
			Embedder:  "mxbai-embed-large",
			Prompt:    "",
			Tokens:    1024,
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

func (o *OllamaAPI) WithTokens(tokens int) {
	o.Tokens = tokens
}

func (o *OllamaAPI) SendMessageTo(ctx context.Context) (*OllamaResponse, error) {
	apiUrl := o.Url + "generate"
	payload := map[string]interface{}{
		"model":   o.Model,
		"prompt":  o.Prompt,
		"stream":  o.Streaming,
		"options": map[string]interface{}{"num_predict": o.Tokens},
	}
	userData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling json: %v", err)
	}
	var resp OllamaResponse

	err = HttpCaller(POST, apiUrl, userData, &resp)
	if err != nil {
		log.Panic(err)
	}

	return &resp, nil
}

func (o *OllamaAPI) GenerateEmbedding(ctx context.Context, content string) (*OllamaEmbeddingResponse, error) {
	apiUrl := o.Url + "embeddings"
	payload := map[string]interface{}{
		"model":  o.Embedder,
		"prompt": strings.Trim(content, "\n"),
	}

	userData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling json: %v", err)
	}
	var resp OllamaEmbeddingResponse

	err = HttpCaller(POST, apiUrl, userData, &resp)
	if err != nil {
		log.Panic(err)
	}

	return &resp, nil
}

func (o *OllamaAPI) PromptFormater(prompt string, values interface{}) {
	var p strings.Builder
	p.WriteString(prompt)
	if values, ok := values.(map[string]string); ok {
		for key, value := range values {
			p.WriteString(fmt.Sprintf("%s:%s", key, value))
		}
	}
	o.Prompt = p.String()
}

func FormatPromptChat(prompt, context string) string {
	fprompt := fmt.Sprintf(`
    You are a helpfull assistant. Please answer the question provided in the PROMPT.
    PROMPT: %s
  `, prompt, context)
	return fprompt
}

func FormatPromptRag(prompt, context string) string {
	fprompt := fmt.Sprintf(`
    You are a helpfull assistant that provides answer based on the knowledge given to you.
    If There is no context, answer: "I don't know, maybe I need more context".

    Use the following context if not blank: %s
    Answer the following:
    PROMPT: %s
  `, prompt, context)
	return fprompt
}
