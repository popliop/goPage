package server

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

// Data for GPT
const systemContent = `Provide the most accurate Swedish export tariff code for the given product. 
-The code must be exactly 6 digits.
- Indicate your confidence level as a percentage.
- Only respond if your confidence is above 95%.
- If unsure or unable to provide an answer, respond with 'Unable to determine'.
- Do not include any additional text or explanation beyond the code and confidence level.`

const gptModel = "gpt-3.5-turbo"

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

func newGPTPost(userContent string) GPTPost {

	jsonBody := fmt.Sprintf(`{
	"model": "%s",
	"messages": [
			{"role": "system", "content": "%s"},
			{"role": "user", "content": "%s"}
		]
	}`, gptModel, strings.ReplaceAll(systemContent, "\n", ""), userContent)

	fmt.Println(jsonBody)

	return GPTPost{
		URL:    "https://api.openai.com/v1/chat/completions",
		Body:   strings.NewReader(jsonBody),
		Method: "POST",
		ApiKey: os.Getenv("GPT_API_KEY"),
	}
}

func sendtoGPT(item string) (responseBack string) {

	gpt := newGPTPost(string(item))

	req, err := http.NewRequest(gpt.Method, gpt.URL, gpt.Body)
	if err != nil {
		log.Fatal(err)
	}

	setHeaders(req, gpt.ApiKey)

	result, err := handleResponse(req)
	if err != nil {
		log.Println("Error handle response")
	}

	if len(result.Choices) > 0 && result.Choices[0].Message.Content != "" {
		return result.Choices[0].Message.Content
	} else {
		log.Println("Choices slice is empty or content is missing")
		return "" // Or handle it appropriately (e.g., return an error)
	}
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

func handleResponse(r *http.Request) (*GPTResponse, error) {
	response, err := sendRequest(r)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	result, err := parseJson(responseBody)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func parseJson(body []byte) (*GPTResponse, error) {
	var result GPTResponse

	err := json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
