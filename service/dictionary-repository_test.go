package service

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_dictionaryRepository_Load(t *testing.T) {

	dict, err := dictRepository.Load("uw", "AbsenceType", "")

	assert.NoError(t, err)

	fmt.Printf("%#v\n", dict)
}

func Test_dictionaryRepository_LoadChildren(t *testing.T) {

	dict, err := dictRepository.LoadChildren("uw", "AbsenceType", "")

	assert.NoError(t, err)

	fmt.Printf("%#v\n", dict)
}

func Test_dictionaryRepository_LoadAll(t *testing.T) {

	dicts, err := dictRepository.LoadAll("")

	assert.NoError(t, err)

	for _, d := range dicts {
		fmt.Printf("%#v\n", d)
	}

}

func Test_dictionaryRepository_LoadByType(t *testing.T) {

}

func Test_dictionaryRepository_Save(t *testing.T) {

}

func Test_dictionaryRepository_LoadChildrenKeys(t *testing.T) {

	res, err := dictRepository.LoadChildrenKeys("uw", "AbsenceType", "")
	assert.NoError(t, err)

	fmt.Printf("%v", res)
}
