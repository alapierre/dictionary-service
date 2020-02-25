package service

import (
	"fmt"
	"testing"
)

func TestDictionaryService_Load(t *testing.T) {
	dict, err := service.Load("uw", "AbsenceType", "")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", dict)
}

func TestDictionaryService_LoadShallow(t *testing.T) {
}

func TestNewDictionaryService(t *testing.T) {

}

func Test_mergeMaps(t *testing.T) {

}

func Test_prepareChildrenMap(t *testing.T) {

}

func Test_prepareMap(t *testing.T) {

}
