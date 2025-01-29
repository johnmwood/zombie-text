package main

import (
	"fmt"
	"os"

	"github.com/johnmwood/zombie-text/internal/analyzer/claude"
	"github.com/johnmwood/zombie-text/internal/config"
)

const (
	secretPath = ".secrets.json"
	configPath = "config.json" // regular config
)

func main() {
	pwd, _ := os.Getwd()
	fmt.Printf("pwd: %s\n", pwd)

	cfg, err := config.LoadConfig(configPath, secretPath)
	if err != nil {
		panic(err)
	}
	analyzer, err := claude.NewAnthropicAnalyzer(cfg)
	if err != nil {
		panic(err)
	}

	err = analyzer.ReadImage()
	fmt.Println(err)
}
