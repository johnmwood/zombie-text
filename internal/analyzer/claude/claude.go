package claude

import (
	"fmt"

	"github.com/johnmwood/zombie-text/internal/config"

	"github.com/liushuangls/go-anthropic/v2"
)

type AnthropicAnalyzer struct {
	Client   anthropic.Client
	ImageDir string
}

func NewAnthropicAnalyzer(cfg *config.Config) (*AnthropicAnalyzer, error) {
	analyzer := &AnthropicAnalyzer{
		Client:   *anthropic.NewClient(cfg.ClaudeAPIKey),
		ImageDir: cfg.ImageDir,
	}
	if cfg.ClaudeAPIKey == "" {
		return nil, fmt.Errorf("missing claude api key")
	}

	return analyzer, nil
}
