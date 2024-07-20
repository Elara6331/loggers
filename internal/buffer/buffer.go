package buffer

import (
	"bytes"
	"sync"
)

const (
	maxBufSize  = 16384
	initBufSize = 1024
)

var bufPool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(make([]byte, 0, initBufSize))
	},
}

// Alloc gets an existing buffer or allocates a new 1kB buffer.
func Alloc() *bytes.Buffer {
	return bufPool.Get().(*bytes.Buffer)
}

// Free releases the buffer back into the pool. If the buffer
// is over 16kB, it isn't put back into the pool to avoid excessive
// memory usage.
func Free(buf *bytes.Buffer) {
	if buf.Cap() <= maxBufSize {
		buf.Reset()
		bufPool.Put(buf)
	}
}
