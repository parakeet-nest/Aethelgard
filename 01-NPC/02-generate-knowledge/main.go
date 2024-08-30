package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/parakeet-nest/parakeet/content"
	"github.com/parakeet-nest/parakeet/embeddings"
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
	var embeddingsModel = "all-minilm:33m" // This model is for the embeddings of the documents 67 MB
	cert, _ := os.ReadFile(os.Getenv("ELASTIC_CERT_PATH"))

	elasticStore := embeddings.ElasticSearchStore{}
	err = elasticStore.Initialize(
		[]string{
			os.Getenv("ELASTIC_ADDRESS"),
		},
		os.Getenv("ELASTIC_USERNAME"),
		os.Getenv("ELASTIC_PASSWORD"),
		cert,
		characterSheetId + "-index",
	)

	if err != nil {
		log.Fatalln("😡:", err)
	}

	characterSheet, err := content.ReadTextFile("./character-sheet-"+characterSheetId+".md")
	if err != nil {
		log.Fatalln("😡:", err)
	}

	//chunks := content.SplitTextWithRegex(characterSheet, `## *`)
	//chunks := content.SplitTextWithDelimiter(characterSheet, "\n\n")
	chunks := content.SplitTextWithRegex(characterSheet, `# *`)

	// Create embeddings from documents and save them in the store
	for idx, doc := range chunks {
		fmt.Println("Creating embedding from document ", idx)
		embedding, err := embeddings.CreateEmbedding(
			ollamaUrl,
			llm.Query4Embedding{
				Model:  embeddingsModel,
				Prompt: doc,
			},
			strconv.Itoa(idx),
		)
		if err != nil {
			fmt.Println("😡:", err)
		} else {
			_, err := elasticStore.Save(embedding)
			if err != nil {
				fmt.Println("😡:", err)
			} else {
				fmt.Println("Document", embedding.Id, "indexed successfully")
			}
		}
	}
}
