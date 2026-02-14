package store

import (
	"context"
	"net/http"
	"net/url"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/vectorstores/qdrant"
)

// We define the struct wrapper to make it easier to reference in the bot package
type QdrantStore struct {
	*qdrant.Store
}

func NewQdrantStore(e embeddings.Embedder) (*qdrant.Store, error) {
	ctx := context.Background()
	u, _ := url.Parse("http://localhost:6333")
	collectionName := "personal-bio"

	// Configuration for the vector space.
	// nomic-embed-text typically uses 768 dimensions.
	collectionConfig := map[string]interface{}{
		"vectors": map[string]interface{}{
			"size":     768,
			"distance": "Cosine",
		},
	}

	// Ensure the collection exists in Qdrant
	urlCollection := u.JoinPath("collections", collectionName)
	// We ignore these returns for now, but in production, you'd check if the collection 409s (already exists)
	_, _, _ = qdrant.DoRequest(ctx, *urlCollection, "", http.MethodPut, collectionConfig)

	s, err := qdrant.New(
		qdrant.WithURL(*u),
		qdrant.WithEmbedder(e),
		qdrant.WithCollectionName(collectionName),
	)

	if err != nil {
		return nil, err
	}

	// FIX: s is a value type, so we take its address to return a pointer.
	return &s, nil
}
