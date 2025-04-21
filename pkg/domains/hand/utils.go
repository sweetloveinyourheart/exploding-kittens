package hand

import (
	"bytes"
	"sync"
)

var bytesBufferPool = sync.Pool{
	New: func() any { return new(bytes.Buffer) },
}
