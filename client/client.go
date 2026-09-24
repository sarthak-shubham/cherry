package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var serverURL = "http://localhost:8080"

var httpClient = &http.Client{}

type setRequest struct {
	Value string `json:"value"`
}

type kvResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func doRequest(method string, key string, body io.Reader) (*http.Response, error) {
	target := fmt.Sprintf("%s/kv/%s", serverURL, url.PathEscape(key))

	request, err := http.NewRequest(method, target, body)
	if err != nil {
		return nil, err
	}

	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	return httpClient.Do(request)
}

func getKey(key string) {
	response, err := doRequest(http.MethodGet, key, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		printAPIError(response)
		return
	}

	var result kvResponse

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s = %s\n", result.Key, result.Value)
}

func setKey(key string, value string) {
	requestBody, err := json.Marshal(setRequest{
		Value: value,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	response, err := doRequest(
		http.MethodPut,
		key,
		bytes.NewReader(requestBody),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		printAPIError(response)
		return
	}

	var result kvResponse

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s = %s\n", result.Key, result.Value)
}

func deleteKey(key string) {
	response, err := doRequest(http.MethodDelete, key, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		printAPIError(response)
		return
	}

	var result kvResponse

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("deleted %s\n", result.Key)
}

func printAPIError(response *http.Response) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("request failed with status %d\n", response.StatusCode)
		return
	}

	var result errorResponse

	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Printf("request failed with status %d\n", response.StatusCode)
		return
	}

	fmt.Printf("error: %s\n", result.Error)
}
