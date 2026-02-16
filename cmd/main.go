package main

import (
	"1_basic_RAG/internal/bot"
	"1_basic_RAG/internal/config"
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

	// 1. Load the environment configuration
	envConfig := config.Load()
	if envConfig.TELEGRAM_APITOKEN == "" {
		log.Fatal("TELEGRAM_APITOKEN is not set in the environment")
	}

	// Get Token from environment variable for security
	token := os.Getenv("TELEGRAM_APITOKEN")
	if token == "" {
		log.Fatal("Please set TELEGRAM_APITOKEN environment variable")
	}

	// 2. Init RAG Components
	llm, _ := ollama.New(ollama.WithModel("llama3"))
	embedClient, _ := ollama.New(ollama.WithModel("nomic-embed-text"))
	embedder, _ := embeddings.NewEmbedder(embedClient)

	// 3. Init Store and Service
	qStore, _ := store.NewQdrantStore(embedder)
	ragService := rag.NewService(llm, qStore)

	// 4. Init and Start Telegram Bot
	telegramBot, err := bot.NewTelegramBot(token, ragService, &store.QdrantStore{Store: qStore})
	if err != nil {
		log.Fatalf("Failed to init bot: %v", err)
	}

	log.Println("Bot is running...")
	telegramBot.Start(ctx)
}
