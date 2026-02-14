package document

import (
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
)

// ProcessText converts raw strings into chunked Documents
func ProcessText(input string) ([]schema.Document, error) {
	// Define how we want to split the local files
	splitter := textsplitter.NewRecursiveCharacter()
	splitter.ChunkSize = 500
	splitter.ChunkOverlap = 50

	return textsplitter.CreateDocuments(splitter, []string{input}, nil)
}
