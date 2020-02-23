package service

import (
	"fmt"
	"testing"
)

func Test_dictionaryRepository_Load(t *testing.T) {

	dict, err := dictRepository.Load("uw", "AbsenceType", "")

	if err != nil {
		t.Fatal(err, "Can't query")
	}

	fmt.Printf("%#v\n", dict)
}

func Test_dictionaryRepository_LoadAll(t *testing.T) {

	dicts, err := dictRepository.LoadAll("")

	if err != nil {
		t.Fatal(err, "Can't query")
	}

	for _, d := range dicts {
		fmt.Printf("%#v\n", d)
	}

}

func Test_dictionaryRepository_LoadByType(t *testing.T) {

}

func Test_dictionaryRepository_Save(t *testing.T) {

}
