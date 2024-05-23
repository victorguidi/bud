package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// HTTPMethod defines a custom type for HTTP methods
type HTTPMethod string

// Constants for common HTTP methods
const (
	GET     HTTPMethod = "GET"
	POST    HTTPMethod = "POST"
	PUT     HTTPMethod = "PUT"
	DELETE  HTTPMethod = "DELETE"
	PATCH   HTTPMethod = "PATCH"
	HEAD    HTTPMethod = "HEAD"
	OPTIONS HTTPMethod = "OPTIONS"
)

func (m HTTPMethod) String() string {
	return string(m)
}

func HttpCaller[T any](method HTTPMethod, apiUrl string, userData []byte, responseBody *T) error {
	// create new http request
	request, err := http.NewRequest(method.String(), apiUrl, bytes.NewBuffer(userData))
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		return err
	}

	// send the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return err
	}
	// I need to encode the response in the resp
	fmt.Println("Status: ", response.Status)

	// clean up memory after execution
	defer response.Body.Close()
	return nil
}
