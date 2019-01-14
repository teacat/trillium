package main

import (
	"fmt"
	"io/ioutil"

	"github.com/teacat/trillium"
)

func main() {
	var mem string

	b := trillium.New(0)

	for i := 0; i < 300000; i++ {
		//<-time.After(time.Millisecond * 1)

		mem = fmt.Sprintf("%s\n%d", mem, b.Generate())
		//fmt.Printf("Its: %d", b.Generate())
	}
	ioutil.WriteFile("yee.txt", []byte(mem), 0777)
}
