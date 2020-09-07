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

func TestTimeout(t *testing.T) {
	assert := assert.New(t)

	var i int
	for {
		if i > 512 {
			break
		}
		i++

		_, err := New(DefaultConfig()).Generate()
		assert.NoError(err)
	}
}

func TestGoroutine(t *testing.T) {
	assert := assert.New(t)
	tri := New(DefaultConfig())

	done := make(chan bool, 2)
	go func() {
		var i int
		for {
			if i > 512 {
				done <- true
				break
			}
			i++

			_, err := tri.Generate()
			assert.NoError(err)
		}
	}()
	go func() {
		var i int
		for {
			if i > 512 {
				done <- true
				break
			}
			i++

			_, err := tri.Generate()
			assert.NoError(err)
		}
	}()
	<-done
}

func BenchmarkUint64(b *testing.B) {
	t := New(DefaultConfig())
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		t.Generate()
	}
}
