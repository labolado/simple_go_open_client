package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Client struct {
	apiURL   string
	apiKey   string
	model    string
	debug    bool
	debugLog *log.Logger
	client   *http.Client
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Choices []struct {
		Index   int     `json:"index"`
		Message Message `json:"message"`
	} `json:"choices"`
}

func NewClient(apiURL, apiKey, model string) *Client {
	return &Client{
		apiURL: apiURL,
		apiKey: apiKey,
		model:  model,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) SetDebug(debug bool) {
	c.debug = debug
	if debug {
		c.debugLog = log.New(os.Stderr, "DEBUG: ", log.Ltime|log.Lshortfile)
	} else {
		c.debugLog = nil
	}
}

func (c *Client) debugPrint(format string, v ...interface{}) {
	if c.debug && c.debugLog != nil {
		c.debugLog.Printf(format, v...)
	}
}
func (c *Client) ChatCompletion(systemContent, userContent string, temperature float64) (string, error) {
    return c.ChatCompletionWithContext(context.Background(), systemContent, userContent, temperature)
	 
}
func (c *Client) ChatCompletionWithContext(ctx context.Context, systemContent, userContent string, temperature float64) (string, error) {
	url := c.apiURL + "/v1/chat/completions"

	messages := []Message{
		{Role: "system", Content: systemContent},
		{Role: "user", Content: userContent},
	}

	requestBody := ChatCompletionRequest{
		Model:       c.model,
		Messages:    messages,
		Temperature: temperature,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshalling request body: %w", err)
	}

	c.debugPrint("Request body: %s", string(jsonBody))

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	c.debugPrint("Response body: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	var chatResponse ChatCompletionResponse
	if err := json.Unmarshal(body, &chatResponse); err != nil {
		return "", fmt.Errorf("error unmarshalling response: %w", err)
	}

	if len(chatResponse.Choices) == 0 {
		return "", fmt.Errorf("no choices in the response")
	}

	return chatResponse.Choices[0].Message.Content, nil
}