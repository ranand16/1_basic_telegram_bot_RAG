package main

import (
	"1_basic_RAG/internal/bot"
	"1_basic_RAG/internal/rag"
	"1_basic_RAG/internal/store"
	"context"
	"log"
	"os"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
)

func main() {
	ctx := context.Background()

	// Get Token from environment variable for security
	token := os.Getenv("TELEGRAM_APITOKEN")
	if token == "" {
		log.Fatal("Please set TELEGRAM_APITOKEN environment variable")
	}

	// 1. Init RAG Components
	llm, _ := ollama.New(ollama.WithModel("llama3"))
	embedClient, _ := ollama.New(ollama.WithModel("nomic-embed-text"))
	embedder, _ := embeddings.NewEmbedder(embedClient)

	// 2. Init Store and Service
	qStore, _ := store.NewQdrantStore(embedder)
	ragService := rag.NewService(llm, qStore)

	// 3. Init and Start Telegram Bot
	telegramBot, err := bot.NewTelegramBot(token, ragService, &store.QdrantStore{Store: qStore})
	if err != nil {
		log.Fatalf("Failed to init bot: %v", err)
	}

	log.Println("Bot is running...")
	telegramBot.Start(ctx)
}
