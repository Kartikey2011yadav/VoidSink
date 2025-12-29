package tests

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/Kartikey2011yadav/voidsink/internal/heffalump"
)

// TestScanHTML verifies that our custom tokenizer correctly splits
// text into words and HTML tags.
func TestScanHTML(t *testing.T) {
	input := `<html> <body> Hello, world! <div class="test">`
	expected := []string{"<html>", "<body>", "Hello,", "world!", "<div class=\"test\">"}

	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(heffalump.ScanHTML)

	var result []string
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("Scanner error: %v", err)
	}

	if len(result) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d: %v", len(expected), len(result), result)
	}

	for i, token := range result {
		if token != expected[i] {
			t.Errorf("Token %d: expected %q, got %q", i, expected[i], token)
		}
	}
}

// TestHeffalump_Integration tests the full lifecycle of the engine:
// Loading a corpus, seeding, and generating text.
func TestHeffalump_Integration(t *testing.T) {
	// 1. Create a temporary corpus file
	content := "The quick brown fox jumps over the lazy dog"
	tmpfile, err := os.CreateTemp("", "corpus.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // Clean up

	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// 2. Initialize Heffalump with the temp file
	h, err := heffalump.New(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to create Heffalump: %v", err)
	}

	// 3. Test Seeding
	w1, w2 := h.Seed()
	if w1 == "" || w2 == "" {
		t.Error("Seed returned empty strings")
	}

	// 4. Test Generation (Next)
	// We know "The quick" -> "brown"
	next := h.Next("The", "quick")
	if next != "brown" {
		t.Errorf("Expected 'brown', got %q", next)
	}

	// Test unknown sequence (should return random or void)
	safeNext := h.Next("unknown", "sequence")
	if safeNext == "" {
		t.Error("Next() returned empty string for unknown sequence")
	}
}
