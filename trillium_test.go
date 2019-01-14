package trillium

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicRoute(t *testing.T) {
	assert := assert.New(t)
	b := New(0)

	a, _ := os.OpenFile("yee.txt", os.O_CREATE|os.O_APPEND, os.ModeAppend)
	for i := 0; i < 300000; i++ {
		//<-time.After(time.Millisecond * 1)

		a.WriteString(strconv.Itoa(b.Generate()))
		a.WriteString("\n")

		//fmt.Println(a.Generate())
	}
	a.Close()
	assert.Equal("a", "v")
}
