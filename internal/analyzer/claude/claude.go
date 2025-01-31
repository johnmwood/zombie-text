package claude

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/johnmwood/zombie-text/internal/config"

	anthropic "github.com/liushuangls/go-anthropic/v2"
)

type AnthropicAnalyzer struct {
	Client     anthropic.Client
	ImageDir   string
	basePrompt *string

	model anthropic.Model
}

func NewAnthropicAnalyzer(cfg *config.Config) (*AnthropicAnalyzer, error) {
	prompt, err := gatherPrompt(cfg.BasePrompt)
	if err != nil {
		fmt.Println("failed to find prompt file at %s", cfg.BasePrompt)
	}
	analyzer := &AnthropicAnalyzer{
		Client:     *anthropic.NewClient(cfg.ClaudeAPIKey),
		ImageDir:   cfg.ImageDir,
		model:      anthropic.ModelClaude3Opus20240229,
		basePrompt: &prompt,
	}
	if cfg.ClaudeAPIKey == "" {
		return nil, fmt.Errorf("missing claude api key")
	}

	return analyzer, nil
}

func (aa *AnthropicAnalyzer) ReadImage(imageName string) error {
	image, err := os.Open(fmt.Sprintf("%s/%s", aa.ImageDir, imageName))
	if err != nil {
		return err
	}
	defer image.Close()
	imageData, err := io.ReadAll(image)
	if err != nil {
		return err
	}
	imageMediaType := "image/png"

	resp, err := aa.Client.CreateMessages(context.Background(), anthropic.MessagesRequest{
		Model: aa.model,
		MultiSystem: []anthropic.MessageSystemPart{
			{
				Type: "text",
				Text: *aa.basePrompt,
			},
		},
		Messages: []anthropic.Message{
			{
				Role: anthropic.RoleUser,
				Content: []anthropic.MessageContent{
					anthropic.NewImageMessageContent(
						anthropic.NewMessageContentSource(
							anthropic.MessagesContentSourceTypeBase64,
							imageMediaType,
							imageData,
						),
					),
				},
			},
		},
		MaxTokens: 1000,
	})
	if err != nil {
		var e *anthropic.APIError
		if errors.As(err, &e) {
			return fmt.Errorf("messages error, type: %q\nmessage: %q\n", e.Type, e.Message)
		}
		return fmt.Errorf("messages error: %v\n", err)
	}
	fmt.Println(resp.Content[0].GetText())

	return nil
}

func gatherPrompt(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
