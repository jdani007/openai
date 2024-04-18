package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

func main() {
	
	context, err := generatePrompt()
	if err != nil {
		log.Fatal(err)
	}

	for {
		response, err := getCompletion(0, context)
		if err != nil {
			log.Fatal(err)
		}

		context = append(context, openai.ChatCompletionMessage{
			Role: openai.ChatMessageRoleAssistant,
			Content: response,
		})

		content, err := getInput()
		if err != nil {
			log.Fatal(err)
		}

		if content == "exit" {
			fmt.Println("\nAssistant: Goodbye!")
			return
		}

		context = append(context, openai.ChatCompletionMessage{
			Role: openai.ChatMessageRoleUser,
			Content: content,
		})
	}
}

func generatePrompt() ([]openai.ChatCompletionMessage, error) {
	content, err := getInput()
	if err != nil {
		return nil, err
	}

	return []openai.ChatCompletionMessage {
		{
			Role: openai.ChatMessageRoleSystem,
			Content: content,
		},
	}, nil
}

func getCompletion(temp float32, ctx []openai.ChatCompletionMessage) (string, error) {
	creds, ok := os.LookupEnv("OPENAPI")
	if !ok {
		return "", fmt.Errorf("missing environment variable 'OPENAPI'")
	}

	client := openai.NewClient(creds)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: ctx,
			Temperature: temp,
		},
	)
	if err != nil {
		return "", fmt.Errorf("ChatCompletion error: %v", err)
	}
	fmt.Println("\nAssistant:",resp.Choices[0].Message.Content)

	return resp.Choices[0].Message.Content, nil
}

func getInput() (string, error) {
	
	fmt.Print("\nUser: ")
	r := bufio.NewReader(os.Stdin)
	s, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(s), nil
}