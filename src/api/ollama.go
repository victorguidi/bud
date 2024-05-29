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
	DEFAULTPROMPT = `
You are a helpful assistant. Please answer the question provided in the "Input".

Input: %s`

	DEFAULTRAGPROMPT = `
You are a helpful assistant that provides answers based on the knowledge given to you. Use the provided context to answer the question. If you cannot find the answer in the provided context, respond with: "I didn't find anything on the documents provided."

Context:
%s

Input: %s

Answer:
`

	DEFAULTCLASSIFIER = `
You are a text classification model. Your task is to classify the given command as one of the following categories: "chat", "rag", or "kill". Here are the definitions for each category:

"chat": Commands that involve general conversation or questions about information, such as asking for facts or details.
"rag": Commands that involve retrieving specific documents or data, such as asking for reports or papers.
"kill": Commands that involve stopping or terminating a process, such as stopping a worker or a service.
Respond only with the appropriate category name.

For example:

Input: What is the capital of New Zealand?
Output: chat

Input: Stop chat
Output: kill

Input: Do I have any paper about bitcoin?
Output: rag

Now, provide the command to be classified:

Input: %s
Output:
  `
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

func (o *OllamaAPI) WithUrl(url string) *OllamaAPI {
	o.Url = url
	return o
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
	if values, ok := values.(map[string]string); ok {
		vs := []any{}
		for k, value := range values {
			if strings.Contains(prompt, k) {
				vs = append(vs, value)
			}
		}
		p.WriteString(fmt.Sprintf(prompt, vs...))
	}
	log.Println("PROMPT SENT TO MODEL: ", p.String())
	o.Prompt = p.String()
}
