package heffalump

import (
	"bytes"
	"sync"
)

// BufferPool manages a pool of byte buffers to reduce garbage collection overhead
// during high-volume text streaming.
type BufferPool struct {
	pool sync.Pool
}

// NewBufferPool creates a new pool for byte buffers.
func NewBufferPool() *BufferPool {
	return &BufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				// Pre-allocate a buffer of 4KB (standard page size)
				// This is usually enough for a chunk of generated text before flushing.
				return bytes.NewBuffer(make([]byte, 0, 4096))
			},
		},
	}
}

// Get retrieves a buffer from the pool.
func (bp *BufferPool) Get() *bytes.Buffer {
	buf := bp.pool.Get().(*bytes.Buffer)
	buf.Reset() // Ensure the buffer is empty before reuse
	return buf
}

// Put returns a buffer to the pool.
func (bp *BufferPool) Put(buf *bytes.Buffer) {
	bp.pool.Put(buf)
}
