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

const testImage = "crossfit_00.PNG"

type AnthropicAnalyzer struct {
	Client   anthropic.Client
	ImageDir string

	model anthropic.Model
}

func NewAnthropicAnalyzer(cfg *config.Config) (*AnthropicAnalyzer, error) {
	analyzer := &AnthropicAnalyzer{
		Client:   *anthropic.NewClient(cfg.ClaudeAPIKey),
		ImageDir: cfg.ImageDir,
		model:    anthropic.ModelClaude3Opus20240229,
	}
	if cfg.ClaudeAPIKey == "" {
		return nil, fmt.Errorf("missing claude api key")
	}

	return analyzer, nil
}

func (aa *AnthropicAnalyzer) ReadImage() error {
	image, err := os.Open(fmt.Sprintf("%s/%s", aa.ImageDir, testImage))
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
