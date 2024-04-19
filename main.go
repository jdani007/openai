package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

const (
	temperature = 0
	greeting    = "Hello"
)

func main() {

	context, client, err := generatePrompt()
	logError(err)

	for {
		summary, err := getCompletion(context, client)
		logError(err)

		fmt.Println("\nAssistant:", summary)

		context = updateContext(summary, openai.ChatMessageRoleAssistant, context)

		userPrompt, err := getInput()
		logError(err)

		switch userPrompt {
		case "generate image":
			fmt.Print("\nAssistant: Generating image")
			go timer()

			err := generateImage(summary, client)
			logError(err)
			return

		case "reset":
			context, _, err = generatePrompt()
			logError(err)

		case "exit", "quit":
			fmt.Print("\nAssistant: Goodbye!\n")
			return

		default:
			context = updateContext(userPrompt, openai.ChatMessageRoleUser, context)
		}
	}

}

func generatePrompt() ([]openai.ChatCompletionMessage, *openai.Client, error) {

	client, err := newClient()
	if err != nil {
		return nil, nil, err
	}

	ctx := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: greeting,
		},
	}

	return ctx, client, nil

}

func getCompletion(ctx []openai.ChatCompletionMessage, client *openai.Client) (string, error) {

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			Messages:    ctx,
			Temperature: temperature,
		},
	)
	if err != nil {
		return "", fmt.Errorf("ChatCompletion error: %v", err)
	}

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

func updateContext(content, role string, ctx []openai.ChatCompletionMessage) []openai.ChatCompletionMessage {

	ctx = append(ctx, openai.ChatCompletionMessage{
		Role:    role,
		Content: content,
	})

	return ctx
}

func timer() {
	for {
		fmt.Print(".")
		time.Sleep(time.Second * 1)
	}
}

func logError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
