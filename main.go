package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

func main() {

	context, err := generatePrompt()
	if err != nil {
		log.Fatal(err)
	}
	client, err := newClient()
	if err != nil {
		log.Fatal(err)
	}

	var summary string
	for {
		summary, err = getCompletion(0, context, client)
		if err != nil {
			log.Fatal(err)
		}

		context = append(context, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: summary,
		})

		userPrompt, err := getInput()
		if err != nil {
			log.Fatal(err)
		}

		if userPrompt == "exit" {
			fmt.Printf("\nAssistant: Goodbye!\n")
			break
		}

		context = append(context, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: userPrompt,
		})
	}
	if err := generateImage(summary, client); err != nil {
		log.Fatal(err)
	}

}

func generatePrompt() ([]openai.ChatCompletionMessage, error) {

	content, err := getInput()
	if err != nil {
		return nil, err
	}

	return []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: content,
		},
	}, nil
}

func getCompletion(temp float32, ctx []openai.ChatCompletionMessage, client *openai.Client) (string, error) {

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			Messages:    ctx,
			Temperature: temp,
		},
	)
	if err != nil {
		return "", fmt.Errorf("ChatCompletion error: %v", err)
	}

	fmt.Println("\nAssistant:", resp.Choices[0].Message.Content)

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

func newClient() (*openai.Client, error) {
	creds, ok := os.LookupEnv("OPENAI")
	if !ok {
		return nil, fmt.Errorf("missing environment variable 'OPENAI'")
	}

	return openai.NewClient(creds), nil
}

func generateImage(summary string, client *openai.Client) error {

	fmt.Println(summary)

	resp, err := client.CreateImage(
		context.Background(),
		openai.ImageRequest{
			Model:   openai.CreateImageModelDallE3,
			Prompt:  summary,
			Size:    openai.CreateImageSize1024x1024,
			Quality: openai.CreateImageQualityStandard,
			N:       1,
		},
	)
	if err != nil {
		return fmt.Errorf("image creation error: %v", err)
	}

	if err := displayImage(resp.Data[0].URL); err != nil {
		return err
	}

	return nil
}

func displayImage(url string) error {

	cmd := exec.Command("open", url)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
