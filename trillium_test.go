package trillium

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUint64(t *testing.T) {
	assert := assert.New(t)
	num, err := New(DefaultConfig()).Generate()
	assert.NoError(err)
	assert.Greater(len(strconv.Itoa(int(num))), 10)
}

func BenchmarkUint64(b *testing.B) {
	t := New(DefaultConfig())
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		t.Generate()
	}
}
