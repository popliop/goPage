package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

type GPTResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type GPTPost struct {
	URL    string
	Body   io.Reader
	Method string
	ApiKey string
}

func newGPTPost() GPTPost {
	return GPTPost{
		URL:    "https://api.openai.com/v1/chat/completions",
		Body:   strings.NewReader(`{"model": "gpt-3.5-turbo", "messages": [{"role": "system", "content": "ONLY GIVE ACCURATE SWEDISH TARRIF CODE. ONLY DIGITS NO FURTHER TEXT"},{"role": "user", "content": "telephone"}]}`),
		Method: "POST",
		ApiKey: os.Getenv("GPT_API_KEY"),
	}
}

func main() {

	gpt := newGPTPost()

	req, err := http.NewRequest(gpt.Method, gpt.URL, gpt.Body)
	if err != nil {
		log.Fatal(err)
	}

	setHeaders(req, newGPTPost().ApiKey)

	response, err := sendRequest(req)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	result, err := parseJson(responseBody)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result.Choices[0].Message.Content)

}

func setHeaders(req *http.Request, apikey string) {
	req.Header.Set("Authorization", "Bearer "+apikey)
	req.Header.Set("Content-Type", "application/json")
}

func sendRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("bad response: %w", err)
	}
	return response, nil
}

func parseJson(body []byte) (*GPTResponse, error) {
	var result GPTResponse

	err := json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
