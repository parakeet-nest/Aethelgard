package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"
)

type Character struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
}

func main() {

	// Read the JSON file
	file, err := os.ReadFile("./character.json")
	if err != nil {
		log.Fatal("😡:", err)
	}

	// Unmarshal the JSON data into a struct
	var character Character
	err = json.Unmarshal(file, &character)
	if err != nil {
		log.Fatal("😡:", err)
	}

	characterSheetId := strings.ToLower(strings.ReplaceAll(character.Name, " ", "-"))
	

	ollamaUrl := "http://localhost:11434"

	//model := "gemma2:2b" // 1.6 GB
	model := "qwen2:1.5b-instruct" // 934 MB
	//model := "phi3:instruct" // 2.2 GB

	instructionsContent, err := os.ReadFile("instructions.md")
	if err != nil {
		log.Fatal("😡:", err)
	}
	stepsContent, err := os.ReadFile("steps.md")
	if err != nil {
		log.Fatal("😡:", err)
	}

	userContent := fmt.Sprintf("Create a %s with this name:%s", character.Kind, character.Name)

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: string(instructionsContent)},
			{Role: "system", Content: string(stepsContent)},
			{Role: "user", Content: userContent},
		},
		Options: llm.Options{
			Temperature:   0.5,
			RepeatLastN:   3,
			RepeatPenalty: 2.0,
			TopK:          10,
			TopP:          0.5,
			//Verbose: true,
		},
	}

	fmt.Println("")
	fmt.Println("🤖 answer:")

	// Answer the question
	answer, err := completion.ChatStream(ollamaUrl, query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		})

	if err != nil {
		log.Fatal("😡:", err)
	}
	// Character sheet
	err = os.WriteFile("../02-generate-knowledge/character-sheet-"+characterSheetId+".md", []byte("# CHARACTER SHEET\n\n"+answer.Message.Content), 0644)
	if err != nil {
		log.Fatal(err)
	}

}
