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

	//model := "qwen2:0.5b" // 352 MB
	model := "phi3:mini" // 2.2 GB

	instructionsContent, err := os.ReadFile("instructions.md")
	if err != nil {
		log.Fatal("😡:", err)
	}
	stepsContent, err := os.ReadFile("steps.md")
	if err != nil {
		log.Fatal("😡:", err)
	}
	userContent := "Give a name for a Dwarf."
	//userContent := "Give a name for an Elf."
	//userContent := "Give a name for a Human."


	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: string(instructionsContent)},
			{Role: "system", Content: string(stepsContent)},
			{Role: "user", Content: userContent},
		},
		Options: llm.Options{
			Temperature:   0.8,
			RepeatLastN:   3,
			RepeatPenalty: 2.0,
			TopK:          10,
			TopP:          0.5,
			Verbose: true,
		},
		Format: "json",
		Raw: true, 
	}

	fmt.Println("")
	fmt.Println("🤖 answer:")

	// Answer the question
	answer, err := completion.Chat(ollamaUrl, query)

	if err != nil {
		log.Fatal("😡:", err)
	}

	fmt.Println(answer.Message.Content)

	err = os.WriteFile("../01-generate-descriptions/character.json", []byte(answer.Message.Content), 0644)
	if err != nil {
		log.Fatal("😡:", err)
	}

	err = os.WriteFile("../02-generate-knowledge/character.json", []byte(answer.Message.Content), 0644)
	if err != nil {
		log.Fatal("😡:", err)
	}
	
	err = os.WriteFile("../03-lets-have-a-chat/character.json", []byte(answer.Message.Content), 0644)
	if err != nil {
		log.Fatal("😡:", err)
	}


}
