//go:build go1.24

package utp

import (
	"runtime"
	"testing"
	"time"
	"weak"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClosedSocketIsCollected(t *testing.T) {
	s, err := NewSocket("udp", "127.0.0.1:0")
	require.NoError(t, err)
	require.NoError(t, s.Close())
	w := weak.Make(s)
	s = nil
	// The packet reader goroutine holds the socket until its read fails.
	for i := 0; i < 100 && w.Value() != nil; i++ {
		runtime.GC()
		time.Sleep(10 * time.Millisecond)
	}
	assert.Nil(t, w.Value())
}
