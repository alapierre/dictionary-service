package service

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"testing"
)

func TestLang(t *testing.T) {

	tag, err := language.Parse("pl")
	assert.NoError(t, err)

	fmt.Printf("%v\n", tag)

	base, c := tag.Base()

	fmt.Printf("%v %v\n", base, c)
}
