package openai

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
)

type Client struct {
    apiURL   string
    apiKey   string
    model    string
    debug    bool
    debugLog *log.Logger
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
        return "", fmt.Errorf("error marshalling request body: %v", err)
    }

    c.debugPrint("Request body: %s", string(jsonBody))

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
    if err != nil {
        return "", fmt.Errorf("error creating request: %v", err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+c.apiKey)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", fmt.Errorf("error sending request: %v", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("error reading response body: %v", err)
    }

    c.debugPrint("Response body: %s", string(body))

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
    }

    var chatResponse ChatCompletionResponse
    err = json.Unmarshal(body, &chatResponse)
    if err != nil {
        return "", fmt.Errorf("error unmarshalling response: %v", err)
    }

    if len(chatResponse.Choices) == 0 {
        return "", fmt.Errorf("no choices in the response")
    }

    return chatResponse.Choices[0].Message.Content, nil
}


