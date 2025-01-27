package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

const (
	systemContent = `Provide the most accurate Swedish export tariff code for the given product. 
-The code must be exactly 6 digits.
- Indicate your confidence level as a percentage.
- Only respond if your confidence is above 95%.
- If unsure or unable to provide an answer, respond with 'Unable to determine'.
- Do not include any additional text or explanation beyond the code and confidence level.`

	gptModel       = "gpt-3.5-turbo"
	gptAPIEndpoint = "https://api.openai.com/v1/chat/completions"
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

type Config struct {
	APIEndpoint string
	Model       string
	APIKey      string
}

type Client struct {
	config     Config
	httpClient *http.Client
}

func NewClient(config Config) *Client {
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) ClassifyProduct(productDescription string) (string, error) {
	request, err := c.createRequest(productDescription)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	response, err := c.executeRequest(request)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}

	if len(response.Choices) == 0 || response.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("invalid response: no content in choices")
	}

	return response.Choices[0].Message.Content, nil
}

func (c *Client) createRequest(userContent string) (*http.Request, error) {
	requestBody := map[string]interface{}{
		"model": c.config.Model,
		"messages": []map[string]string{
			{"role": "system", "content": strings.ReplaceAll(systemContent, "\n", "")},
			{"role": "user", "content": userContent},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", c.config.APIEndpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c *Client) executeRequest(req *http.Request) (*GPTResponse, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	var gptResponse GPTResponse
	if err := json.NewDecoder(resp.Body).Decode(&gptResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &gptResponse, nil
}

func NewGPTPost(userContent string) (*GPTPost, error) {
	requestBody := map[string]interface{}{
		"model": gptModel,
		"messages": []map[string]string{
			{"role": "system", "content": strings.ReplaceAll(systemContent, "\n", "")},
			{"role": "user", "content": userContent},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	return &GPTPost{
		URL:    gptAPIEndpoint,
		Body:   bytes.NewReader(jsonBody),
		Method: "POST",
		ApiKey: os.Getenv("GPT_API_KEY"),
	}, nil
}

func SendToGPT(userInput string) (string, error) {
	gptPost, err := NewGPTPost(userInput)
	if err != nil {
		return "", fmt.Errorf("failed to create GPT request: %w", err)
	}

	req, err := http.NewRequest(gptPost.Method, gptPost.URL, gptPost.Body)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	setHeaders(req, gptPost.ApiKey)

	response, err := executeRequest(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute GPT request: %w", err)
	}

	if len(response.Choices) == 0 || response.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("invalid response from GPT: no content in choices")
	}

	return response.Choices[0].Message.Content, nil
}

func setHeaders(req *http.Request, apiKey string) {
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
}

func executeRequest(req *http.Request) (*GPTResponse, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var gptResponse GPTResponse
	if err := json.Unmarshal(body, &gptResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &gptResponse, nil
}
