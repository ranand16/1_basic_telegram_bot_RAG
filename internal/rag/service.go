package rag

import (
	"context"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/vectorstores"
)

type Service struct {
	llm   llms.Model
	store vectorstores.VectorStore
}

// To cerate new RAG a cordinator
func NewService(l llms.Model, s vectorstores.VectorStore) *Service {
	return &Service{
		llm:   l,
		store: s,
	}
}

// give context and question to get the answer in return
func (s *Service) Ask(ctx context.Context, question string) (string, error) {
	retriever := vectorstores.ToRetriever(s.store, 3)

	chain := chains.NewRetrievalQAFromLLM(s.llm, retriever)

	return chains.Run(ctx, chain, question)
}
