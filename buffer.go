package utils

import "sync"

const (
	bytesBufferSize = 256
)

var bytesBufferPool = sync.Pool{
	New: func() any {
		return &BytesBuffer{buf: make([]byte, bytesBufferSize)}
	},
}

func GetBytesBuffer() *BytesBuffer {
	return bytesBufferPool.Get().(*BytesBuffer)
}

func PutBytesBuffer(buf *BytesBuffer) {
	buf.Reset()
	bytesBufferPool.Put(buf)
}

type BytesBuffer struct {
	buf []byte
	pos int
}

func (b *BytesBuffer) Write(vals []byte) {
	if b.pos+len(vals) > len(b.buf) {
		newBuf := make([]byte, max(2*len(b.buf), b.pos+len(vals)))
		copy(newBuf, b.buf[:b.pos])
		b.buf = newBuf
	}
	copy(b.buf[b.pos:], vals)
	b.pos += len(vals)
}

func (b *BytesBuffer) WriteByte(val byte) error {
	if b.pos+1 > len(b.buf) {
		newBuf := make([]byte, max(2*len(b.buf), b.pos+1))
		copy(newBuf, b.buf[:b.pos])
		b.buf = newBuf
	}
	b.buf[b.pos] = val
	b.pos++
	return nil
}

func (b *BytesBuffer) Reset() {
	b.pos = 0
}

func (b *BytesBuffer) Bytes() []byte {
	return b.buf[:b.pos]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
