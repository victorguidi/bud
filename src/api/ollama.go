package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
	apiUrl := "http://localhost:11434/api/generate"
	userData := []byte(`{"model":"gemma","prompt":"` + msg + `", "stream":false}`)

	// create new http request
	request, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(userData))
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		return nil, err
	}

	// send the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	type Response struct {
		Model      string `json:"model"`
		Created_at string `json:"created_at"`
		Response   string `json:"response"`
	}
	var resp Response
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	// I need to encode the response in the resp
	fmt.Println("Status: ", response.Status)

	// clean up memory after execution
	defer response.Body.Close()
	return resp, nil
}
