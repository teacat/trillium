package trillium

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	assert := assert.New(t)
	tlm := New(DefaultConfig())
	assert.Len(tlm.Generate().String(), 20)
}

func TestInt(t *testing.T) {
	assert := assert.New(t)
	tlm := New(DefaultConfig())
	now := time.Since(time.Unix(900288000, 0))
	assert.Len(tlm.Generate().Int(), len(strconv.Itoa(int(now.Seconds())))+10)
}

func BenchmarkString(b *testing.B) {
	t := New(DefaultConfig())
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = t.Generate().String()
	}
}

func BenchmarkInt(b *testing.B) {
	t := New(DefaultConfig())
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = t.Generate().Int()
	}
}
