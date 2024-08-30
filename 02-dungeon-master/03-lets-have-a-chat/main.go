package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/gear"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/tools"
	"github.com/parakeet-nest/parakeet/wasm"
)

/*
	{
	  "name": "move",
	  "arguments": {
	    "from": 5,
	    "to": 6
	  }
	}
*/
type MoveTool struct {
	Name      string `json:"name"`
	Arguments struct {
		From int `json:"from"`
		To   int `json:"to"`
	} `json:"arguments"`
}

func parseForMoveQuestion(ollamaUrl, model, question string) (string, error) {

	systemContentIntroduction := `You have access to the following tools:`

	toolsList := []llm.Tool{
		{
			Type: "function",
			Function: llm.Function{
				Name:        "move",
				Description: "I want to move or go from a given room to another with given room ids. Use only the move or go command.",
				Parameters: llm.Parameters{
					Type: "object",
					Properties: map[string]llm.Property{
						"from": {
							Type:        "number",
							Description: "current room id",
						},
						"to": {
							Type:        "number",
							Description: "next room id",
						},
					},
					Required: []string{"from", "to"},
				},
			},
		},
	}
	toolsContent, err := tools.GenerateContent(toolsList)

	/*
	[AVAILABLE_TOOLS] [
		{
			"type":"function",
			"function":{
				"name":"move",
				"description":"I want to move or go from a given room to another with given room ids. Use only the move or go command.",
				"parameters":{
					"type":"object",
					"properties":{"from":{"type":"number","description":"current room id"},"to":{"type":"number","description":"next room id"}},
					"required":["from","to"]
				}
			}
		}
	] [/AVAILABLE_TOOLS]
	*/


	if err != nil {
		return "", err
	}

	// Use this if the LLM do not implement the tools natively
	systemContentInstructions := `If the question of the user matched the description of a tool, the tool will be called.
	To call a tool, respond with a JSON object with the following structure: 
	{
	  "name": <name of the called tool>,
	  "arguments": {
	    <name of the argument>: <value of the argument>
	  }
	}
	
	search the name of the tool in the list of tools with the Name field
	`
	options := llm.Options{
		Temperature:   0.0,
		RepeatLastN:   2,
		RepeatPenalty: 2.0,
		Seed:          123,
	}

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: systemContentIntroduction},
			{Role: "system", Content: toolsContent},
			{Role: "system", Content: systemContentInstructions},
			{Role: "user", Content: question},
		},
		Options: options,
		Format:  "json",
		Raw:     true, // try with false
	}

	answer, err := completion.Chat(ollamaUrl, query)
	if err != nil {
		return "", err
	}

	result, err := gear.PrettyString(answer.Message.Content)
	if err != nil {
		return "", err
	}

	return result, nil
}

/*
- parseForMoveQuestion: Ask the model which tool to use -> JSON string
- use the JSON string to instantiate a MoveTool struct
- if OK -> true, MoveTool{Name: "move", Arguments: {From: 5, To: 6}}
- if not OK -> false, empty MoveTool{}
*/
func DoesThePlayerWantToMove(ollamaUrl, model, question string) (bool, MoveTool, error) {

	// Ask the model which tool to use
	checkIfThePlayerWantsToMove, err := parseForMoveQuestion(ollamaUrl, model, question)
	// it should be a JSON string like this:
	// {"name":"move","arguments":{"from":5,"to":6}}
	if err != nil {
		return false, MoveTool{}, err
	}

	moveTool := MoveTool{}

	err = json.Unmarshal([]byte(checkIfThePlayerWantsToMove), &moveTool)
	if err != nil {
		// that means the JSON string is not a MoveTool or it is not a valid JSON string
		// return false / empty MoveTool but not an error
		// then the chat model will handle the question
		return false, MoveTool{}, nil
	}
	if moveTool.Name == "move" {
		return true, moveTool, nil
	} else {
		return false, MoveTool{}, nil
	}
}

func main() {

	ollamaUrl := "http://localhost:11434"
	var smallChatModel = "qwen2:0.5b"  // This model is for the chat completion 352 MB
	var toolModel = "dolphin-phi:2.7b" // This model is for the tool calling 1.6 GB

	// Load the wasm plugin
	wasmPlugin, err := wasm.NewPlugin("./wasm-rust/target/wasm32-unknown-unknown/debug/wasm_rust.wasm", nil)
	if err != nil {
		log.Fatal("🟪 wasm 😡:", err)
	}

	systemContent := `
		You are a Dungeon Master expert at interpreting and answering questions based on provided sources.
		Using only the provided context, answer the user's question
		to the best of your ability using only the resources provided.
		Be verbose!`
	

	// Chat with the Dungeon Master
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("🤖 [🗣️ %s|🛠️ %s] ask me something> ", smallChatModel, toolModel)
		question, _ := reader.ReadString('\n')
		question = strings.TrimSpace(question)

		if question == "bye" {
			break
		}

		// Detect if the player wants to move
		isItMove, result, err := DoesThePlayerWantToMove(ollamaUrl, toolModel, question)
		/*
		- parseForMoveQuestion: Ask the model which tool to use -> JSON string
		- use the JSON string to instantiate a MoveTool struct
		- if OK -> true, MoveTool{Name: "move", Arguments: {From: 5, To: 6}}
		- if not OK -> false, empty MoveTool{}
		*/

		if err != nil {
			log.Fatal("😡:", err)
		}


		if isItMove {
			fmt.Println()
			fmt.Printf("🛠️ %s > action: %s %d => %d", toolModel, result.Name, result.Arguments.From, result.Arguments.To)
			fmt.Println()

			// Convert the MoveTool Arguments struct to a JSON string
			// This JSON string will be used as an argument for the wasm plugin
			jsonData, err := json.Marshal(result.Arguments)
			if err != nil {
				log.Fatal("😡:", err)
			}
			
			// call the function of the wasm plugin
			res, err := wasmPlugin.Call("move_person", jsonData)

			if err != nil {
				log.Fatal("🟪 wasm 😡:", err)
			}
			fmt.Println("🟣 calling WASM function...")

			// display the result
			fmt.Println("💜 result: ", string(res))

			fmt.Println()

		} else { // Let's chat about something else
			fmt.Println()
			fmt.Printf("🗣️ %s > ", smallChatModel)

			queryChat := llm.Query{
				Model: smallChatModel,
				Messages: []llm.Message{
					{Role: "system", Content: systemContent},
					{Role: "user", Content: question},
				},
				Options: llm.Options{
					Temperature:   0.0,
					RepeatLastN:   2,
					RepeatPenalty: 3.0,
					TopK:          10,
					TopP:          0.5,
				},
			}

			// Answer the question
			_, err := completion.ChatStream(ollamaUrl, queryChat,
				func(answer llm.Answer) error {
					fmt.Print(answer.Message.Content)
					return nil
				})

			fmt.Println()
			fmt.Println()
			if err != nil {
				log.Fatal("😡:", err)
			}
		}

	}

}
