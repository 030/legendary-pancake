package main

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func helperReadFile(f string) []byte {
	b, err := ioutil.ReadFile("test/data/" + f + ".yaml")
	if err != nil {
		fmt.Print(err)
	}
	return b
}
func TestSomething(t *testing.T) {
	key := "matchLabels"
	regex := `\s+-\s` + key + `:\n\s+[a-z0-9\/\.\":\s-]+$`
	s, _ := sortElementsInFileByKey(key, regex)
	writeToFile(key, s)

	exp := helperReadFile("matchLabels/expected")
	act := helperReadFile("matchLabels/actual")

	assert.Equal(t, exp, act)
}

func TestSomethingElse(t *testing.T) {
	bla()

	exp := helperReadFile("resources/expected")
	act := helperReadFile("resources/actual")

	assert.Equal(t, exp, act)
}
