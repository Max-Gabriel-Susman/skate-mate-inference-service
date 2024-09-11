package message

import (
	"context"
	"fmt"
	"regexp"

	openai "github.com/sashabaranov/go-openai"
)

type Message struct {
	Text string
}

func NewMessage(text string) Message {
	return Message{Text: text}
}

func (m Message) RespondToMessage(openAIClient *openai.Client) (*openai.ChatCompletionResponse, error) {
	// fmt.Println("message is: ", message) // delete l8r
	fmt.Println("pre parsed message: ", m.Text) // delete l8r

	// Define the regular expression pattern to match the message
	re := regexp.MustCompile(`\[::1\]:\d+: (.+)`)

	// Find the match
	match := re.FindStringSubmatch(m.Text)

	// Check if a match is found
	if len(match) > 1 {
		m.Text = match[1]
		fmt.Println("Parsed message:", m.Text, ":") // delete l8r
	} else {
		fmt.Println("No match found")
	}

	resp, err := openAIClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: m.Text,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return nil, err
	}

	fmt.Println(resp.Choices[0].Message.Content)

	return &resp, nil
}
