// main.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type ChatCompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func createChatCompletion(prompt string) (string, error) {
	url := fmt.Sprintf("%s/openai/deployments/gpt-4o/chat/completions?api-version=2023-05-15", azureopenai.AzureOpenAIEndpoint)

	requestBody := ChatCompletionRequest{
		Model: "gpt-4o",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", azureopenai.AzureOpenAIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("non-OK HTTP status: %s, body: %s", resp.Status, string(bodyBytes))
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var completionResponse ChatCompletionResponse
	err = json.Unmarshal(bodyBytes, &completionResponse)
	if err != nil {
		return "", err
	}

	if len(completionResponse.Choices) > 0 {
		return completionResponse.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no choices returned in response")
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s \"Your prompt here\"", os.Args[0])
	}

	prompt := os.Args[1]
	fmt.Printf("Prompt: %s\n", prompt)

	response, err := createChatCompletion(prompt)
	if err != nil {
		log.Fatalf("Error creating chat completion: %v", err)
	}

	fmt.Printf("AI Response: %s\n", response)
}
