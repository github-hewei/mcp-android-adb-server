package vision

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type Model struct {
	client    *openai.Client
	baseURL   string
	apiKey    string
	modelName string
}

// NewModel creates a new Model instance.
func NewModel(apiKey, modelName, baseUrl string) *Model {
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = baseUrl
	client := openai.NewClientWithConfig(config)

	return &Model{
		client:    client,
		baseURL:   baseUrl,
		apiKey:    apiKey,
		modelName: modelName,
	}
}

// ScreenshotDescription takes a screenshot of the device and returns a description of the image.
func (m *Model) ScreenshotDescription(filename string) (string, error) {
	imageBase64, err := Image2Base64(filename)
	if err != nil {
		return "", fmt.Errorf("failed to convert image to base64: %w", err)
	}

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: strings.ReplaceAll(SystemPrompt, "^", "`"),
		},
		{
			Role: openai.ChatMessageRoleUser,
			MultiContent: []openai.ChatMessagePart{
				{
					Type: openai.ChatMessagePartTypeImageURL,
					ImageURL: &openai.ChatMessageImageURL{
						URL:    fmt.Sprintf("data:image/png;base64,%s", imageBase64),
						Detail: openai.ImageURLDetailAuto,
					},
				},
			},
		},
	}

	resp, err := m.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    m.modelName,
			Messages: messages,
		})

	if err != nil {
		return "", fmt.Errorf("dialog error：%w", err)
	}

	content := resp.Choices[0].Message.Content
	slog.Info("screenshot_description", "file", filename, "response", content)

	return content, nil
}

// Image2Base64 converts an image file to base64 format.
func Image2Base64(filename string) (string, error) {
	imageBytes, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return base64.StdEncoding.EncodeToString(imageBytes), nil
}
