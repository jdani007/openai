package main

import (
	"bufio"
	ctx "context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

const (
	temperature        = 0
	defaultInstruction = "Hello, you are a friendly assistant that responds in a prefessional, tweet friendly manner."
)

func main() {

	context, client, err := generatePrompt()
	if err != nil {
		log.Fatal(err)
	}

	if err := run(context, client); err != nil {
		log.Fatal(err)
	}

	fmt.Println()
}

func run(context []openai.ChatCompletionMessage, client *openai.Client) error {
	for {
		summary, err := getCompletion(context, client)
		if err != nil {
			return err
		}

		fmt.Println("\nAssistant:", summary)
		context = updateContext(summary, openai.ChatMessageRoleAssistant, context)

		userPrompt, err := getInput()
		if err != nil {
			return err
		}

		switch userPrompt {
		case "generate image":
			return generateImage(summary, client)

		case "reset":
			context, _, err = generatePrompt()
			if err != nil {
				return err
			}

		case "exit", "quit":
			fmt.Print("\nAssistant: Goodbye!")
			return nil

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

	context := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: defaultInstruction,
		},
	}

	return context, client, nil

}

func getCompletion(context []openai.ChatCompletionMessage, client *openai.Client) (string, error) {

	resp, err := client.CreateChatCompletion(
		ctx.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			Messages:    context,
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
	c, ok := os.LookupEnv("OPENAI")
	if !ok {
		return nil, fmt.Errorf("missing environment variable 'OPENAI'")
	}

	return openai.NewClient(c), nil
}

func generateImage(summary string, client *openai.Client) error {

	fmt.Print("\nAssistant: Generating image")
	go func() {
		for {
			fmt.Print(".")
			time.Sleep(time.Second * 1)
		}
	}()

	resp, err := client.CreateImage(
		ctx.Background(),
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

func updateContext(content, role string, context []openai.ChatCompletionMessage) []openai.ChatCompletionMessage {

	context = append(context, openai.ChatCompletionMessage{
		Role:    role,
		Content: content,
	})

	return context
}
