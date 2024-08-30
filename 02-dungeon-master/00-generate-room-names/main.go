package main

import (
	"fmt"
	"log"
	"os"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"
)

func main() {
	ollamaUrl := "http://localhost:11434"

	//model := "phi3:mini"
	model := "gemma2:2b" // 1.6 GB 
	//model := "gemma:2b" // 1.7 GB
/*
tinydolphin:latest              0f9dd11f824c    636 MB  2 weeks ago 
tinyllama:latest                2644915ede35    637 MB  2 weeks ago 
*/

	instructionsContent, err := os.ReadFile("instructions.md")
	if err != nil {
		log.Fatal("😡:", err)
	}
	stepsContent, err := os.ReadFile("steps.md")
	if err != nil {
		log.Fatal("😡:", err)
	}
	userContent, err := os.ReadFile("tasks.md")
	if err != nil {
		log.Fatal("😡:", err)
	}

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: string(instructionsContent)},
			{Role: "system", Content: string(stepsContent)},
			{Role: "user", Content: string(userContent)},
		},
		Options: llm.Options{
			Temperature:   1.3,
			RepeatLastN:   3,
			RepeatPenalty: 2.0,
			TopK:          10,
			TopP:          0.5,
			Verbose:       true,
			Stop: []string{`"number": 26`},
		},
		Format: "json",
		Raw:    false,
	}

	fmt.Println("")
	fmt.Println("🤖 answer:")

	// Answer the question
	// Answer the question
	finalAnswer, err := completion.ChatStream(ollamaUrl, query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		})

	if err != nil {
		log.Fatal("😡:", err)
	}

	err = os.WriteFile("./rooms.json", []byte(finalAnswer.Message.Content), 0644)
	if err != nil {
		log.Fatal("😡:", err)
	}

	err = os.WriteFile("../01-generate-descriptions/rooms.json", []byte(finalAnswer.Message.Content), 0644)
	if err != nil {
		log.Fatal("😡:", err)
	}
}
