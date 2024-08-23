package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/parakeet-nest/parakeet/completion"
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
	var embeddingsModel = "all-minilm:33m" // This model is for the embeddings of the documents
	var smallChatModel = "qwen2:0.5b"      // This model is for the chat completion

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

	systemContentTpl := `You are a %s, your name is %s,
	expert at interpreting and answering questions based on provided sources.
	Using only the provided context, answer the user's question 
	to the best of your ability using only the resources provided. 
	Be verbose!`

	systemContent := fmt.Sprintf(systemContentTpl, character.Kind, character.Name)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("🤖 [%s] ask me something> ", smallChatModel)
		question, _ := reader.ReadString('\n')
		question = strings.TrimSpace(question)

		if question == "bye" {
			break
		}

		// Create an embedding from the question
		embeddingFromQuestion, err := embeddings.CreateEmbedding(
			ollamaUrl,
			llm.Query4Embedding{
				Model:  embeddingsModel,
				Prompt: question,
			},
			"question",
		)
		if err != nil {
			log.Fatalln("😡:", err)
		}
		fmt.Println("🔎 searching for similarity...")

		similarities, err := elasticStore.SearchTopNSimilarities(embeddingFromQuestion, 2)

		for _, similarity := range similarities {
			fmt.Println("📝 doc:", similarity.Id, "score:", similarity.Score)
		}

		if err != nil {
			log.Fatalln("😡:", err)
		}

		documentsContent := embeddings.GenerateContentFromSimilarities(similarities)

		queryChat := llm.Query{
			Model: smallChatModel,
			Messages: []llm.Message{
				{Role: "system", Content: systemContent},
				{Role: "system", Content: documentsContent},
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

		fmt.Println()
		fmt.Println("🤖 answer:")

		// Answer the question
		_, err = completion.ChatStream(ollamaUrl, queryChat,
			func(answer llm.Answer) error {
				fmt.Print(answer.Message.Content)
				return nil
			})

		if err != nil {
			log.Fatal("😡:", err)
		}

		fmt.Println()
	}

}
