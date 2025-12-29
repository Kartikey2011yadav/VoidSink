package heffalump

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// Heffalump is a Markov Chain generator designed for high-performance text streaming.
type Heffalump struct {
	chain map[string][]string // Key: "word1 word2", Value: ["word3", "word3", ...]
	keys  []string            // Cache of keys for random starting points
	mu    sync.RWMutex        // Mutex for thread-safe access (if we add dynamic learning later)
	rnd   *rand.Rand
}

// New creates a new Heffalump instance from a source text file.
func New(path string) (*Heffalump, error) {
	h := &Heffalump{
		chain: make(map[string][]string),
		rnd:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	if err := h.load(path); err != nil {
		return nil, err
	}

	return h, nil
}

// load reads the file and builds the Markov chain.
func (h *Heffalump) load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	var w1, w2 string

	// Initialize first two words
	if scanner.Scan() {
		w1 = scanner.Text()
	}
	if scanner.Scan() {
		w2 = scanner.Text()
	}

	// Build the chain
	for scanner.Scan() {
		w3 := scanner.Text()
		key := w1 + " " + w2

		h.chain[key] = append(h.chain[key], w3)

		// Shift window
		w1, w2 = w2, w3
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Cache keys for random access
	h.keys = make([]string, 0, len(h.chain))
	for k := range h.chain {
		h.keys = append(h.keys, k)
	}

	log.Info().Int("triplets", len(h.chain)).Msg("Heffalump Markov chain loaded")
	return nil
}

// Next returns the next probable word based on the previous two words.
// If the sequence (w1, w2) is unknown or a dead end, it picks a random starting point.
func (h *Heffalump) Next(w1, w2 string) string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	key := w1 + " " + w2
	candidates, ok := h.chain[key]

	if !ok || len(candidates) == 0 {
		// Dead end or unknown sequence, pick a random key to restart the flow
		if len(h.keys) == 0 {
			return "void"
		}
		randomKey := h.keys[h.rnd.Intn(len(h.keys))]
		// To restart smoothly, we pretend we just saw the random key.
		// We return the *first* word of that key to reset the caller's state eventually,
		// or we can just return a random word from that key's candidates.
		// Let's return a random candidate from a random key to ensure we always return a "next" word.
		candidates = h.chain[randomKey]
	}

	return candidates[h.rnd.Intn(len(candidates))]
}

// Seed returns a random pair of words to start the generation loop.
func (h *Heffalump) Seed() (string, string) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.keys) == 0 {
		return "the", "void"
	}

	key := h.keys[h.rnd.Intn(len(h.keys))]
	parts := strings.Split(key, " ")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "the", "void"
}
