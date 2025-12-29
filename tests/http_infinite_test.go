package tests

import (
	"bufio"
	"context"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Kartikey2011yadav/voidsink/internal/heffalump"
	"github.com/Kartikey2011yadav/voidsink/internal/trap"
)

func TestHTTPInfiniteTrap(t *testing.T) {
	// 1. Setup Heffalump with dummy data
	tmpfile, err := os.CreateTemp("", "corpus.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.WriteString("word1 word2 word3 word4 word5 word6")
	tmpfile.Close()

	h, err := heffalump.New(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// 2. Create Trap
	// Use a random high port to avoid conflicts
	addr := "127.0.0.1:54321"
	serverName := "fake-nginx"
	tr := trap.NewHTTPInfiniteTrap(addr, serverName, h)

	// 3. Start Trap in a goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		errChan <- tr.Start(ctx)
	}()

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)

	// 4. Make a request
	resp, err := http.Get("http://" + addr + "/test")
	if err != nil {
		t.Fatalf("Failed to connect to trap: %v", err)
	}
	defer resp.Body.Close()

	// 5. Verify Headers
	if got := resp.Header.Get("Server"); got != serverName {
		t.Errorf("Expected Server header %q, got %q", serverName, got)
	}
	if got := resp.Header.Get("Content-Type"); got != "text/html" {
		t.Errorf("Expected Content-Type text/html, got %q", got)
	}

	// 6. Verify Streaming Data
	// Read a small chunk to ensure it's sending data
	reader := bufio.NewReader(resp.Body)
	chunk := make([]byte, 1024)
	n, err := reader.Read(chunk)
	if err != nil {
		t.Fatalf("Failed to read from trap: %v", err)
	}
	if n == 0 {
		t.Error("Trap sent 0 bytes")
	}

	// Verify content looks like words (basic check)
	content := string(chunk[:n])
	if !strings.Contains(content, "word") {
		t.Logf("Warning: Generated content didn't contain expected words. Got: %q", content)
	}

	// 7. Shutdown
	cancel()
	// Wait for shutdown to complete (optional, but good practice)
	select {
	case err := <-errChan:
		if err != nil {
			t.Logf("Trap shutdown error (expected): %v", err)
		}
	case <-time.After(1 * time.Second):
		// Timeout
	}
}
