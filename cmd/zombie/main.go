package main

import (
	"fmt"
	"os"
	"strings"

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

	dir, err := os.ReadDir(analyzer.ImageDir)
	if err != nil {
		fmt.Printf("failed to read dir with err: %v", err)
	}

	for _, entry := range dir {
		name := entry.Name()
		fmt.Println("processing file: ", name)
		if !strings.HasSuffix(strings.ToLower(name), "png") {
			continue
		}
		err := analyzer.ReadImage(name)
		if err != nil {
			fmt.Println(err)
		}
	}

}
