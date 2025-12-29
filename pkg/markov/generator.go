package markov

import (
	"math/rand"
)

// Generator is a placeholder for a Markov Chain text generator.
type Generator struct {
	// In a real implementation, this would hold the chain data.
}

// NewGenerator creates a new Markov generator.
func NewGenerator() *Generator {
	return &Generator{}
}

// Generate returns a chunk of generated text.
// For now, it returns random garbage to simulate the "HellPot" behavior.
func (g *Generator) Generate(size int) []byte {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 "
	b := make([]byte, size)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return b
}
