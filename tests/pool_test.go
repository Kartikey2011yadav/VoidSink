package tests

import (
	"testing"

	"github.com/Kartikey2011yadav/voidsink/internal/heffalump"
)

func TestBufferPool(t *testing.T) {
	pool := heffalump.NewBufferPool()

	// 1. Get a buffer
	buf1 := pool.Get()
	if buf1 == nil {
		t.Fatal("Pool returned nil buffer")
	}
	if buf1.Cap() < 4096 {
		t.Errorf("Expected buffer capacity >= 4096, got %d", buf1.Cap())
	}

	// 2. Write to it
	buf1.WriteString("dirty data")
	if buf1.Len() == 0 {
		t.Error("Buffer should have data")
	}

	// 3. Return it to the pool
	pool.Put(buf1)

	// 4. Get it (or another one) back
	buf2 := pool.Get()

	// 5. Verify it is reset (empty)
	if buf2.Len() != 0 {
		t.Errorf("Buffer from pool was not reset! Len: %d", buf2.Len())
	}

	// 6. Verify capacity is still preserved (optimization check)
	if buf2.Cap() < 4096 {
		t.Errorf("Buffer lost capacity! Cap: %d", buf2.Cap())
	}
}
