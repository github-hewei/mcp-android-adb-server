package vision

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"image"
	"image/color"
	"image/draw"
	"image/png"
	"strconv"

	"github.com/sashabaranov/go-openai"
)

type Model struct {
	client    *openai.Client
	baseURL   string
	apiKey    string
	modelName string
}

// Element represents an element in the screenshot.
type Element struct {
	Type         string `json:"type"`
	Position     string `json:"position"`
	Description  string `json:"description"`
	PositionDesc string `json:"position_desc"`
	Text         string `json:"text"`
	Icon         string `json:"icon"`
	Event        string `json:"event"`
}

// Output represents the output of the model.
type Output struct {
	Description string    `json:"description"`
	Elements    []Element `json:"elements"`
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
						Detail: openai.ImageURLDetailHigh,
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
	go func() {
		_ = DrawElements(content, filename)
	}()

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

// DrawElements draws elements on the image.
func DrawElements(content, filename string) error {
	defer func() {
		if r := recover(); r != nil {
			slog.Info("Recovered from panic:", "err", r)
		}
	}()

	if content == "" {
		return fmt.Errorf("content is empty")
	}

	output := &Output{}
	if err := json.Unmarshal([]byte(content), output); err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}

	// Open the original image
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	// Decode PNG image
	img, err := png.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	// Create a new RGBA image
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	// Draw red boxes for each element
	for _, elem := range output.Elements {
		if elem.Position == "" {
			continue
		}

		// Parse coordinates
		coords := strings.Split(elem.Position, ",")
		if len(coords) != 2 {
			continue
		}

		x, err := strconv.Atoi(strings.TrimSpace(coords[0]))
		if err != nil {
			continue
		}

		y, err := strconv.Atoi(strings.TrimSpace(coords[1]))
		if err != nil {
			continue
		}

		// Draw a 20x20 red box
		red := color.RGBA{R: 255, A: 255}
		size := 20

		// Draw four sides of the box
		for i := -size / 2; i <= size/2; i++ {
			// Draw horizontal lines
			rgba.Set(x+i, y-size/2, red)
			rgba.Set(x+i, y+size/2, red)
			// Draw vertical lines
			rgba.Set(x-size/2, y+i, red)
			rgba.Set(x+size/2, y+i, red)
		}
	}

	// Save the modified image
	outFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	if err := png.Encode(outFile, rgba); err != nil {
		return fmt.Errorf("failed to encode output image: %w", err)
	}

	return nil
}
