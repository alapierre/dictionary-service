package service

import (
	"dictionaries-service/model"
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
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

func TestDictionaryService_SaveParent(t *testing.T) {

	content := make(childrenMap)
	content["lp"] = 1
	content["label"] = "testowy"

	p := model.ParentDictionary{
		Key:     "test",
		Type:    "TestType",
		Name:    "Testowy",
		GroupId: nil,
		Tenant:  "",
		Content: nil,
		Children: []model.ChildDictionary{
			{Key: "testCh1", Name: "dziecko 1", Content: content},
			{Key: "testCh2", Name: "dziecko 2", Content: content},
		},
	}

	if err := service.SaveParent(&p); err != nil {
		t.Errorf("Problem with save parent %v", err)
	}
}

func TestDictionaryService_UpdateParent(t *testing.T) {

	content := make(childrenMap)
	content["lp"] = 1
	content["label"] = "testowy"

	p := model.ParentDictionary{
		Key:     "test",
		Type:    "TestType",
		Name:    "Testowy updated",
		GroupId: nil,
		Tenant:  "",
		Content: nil,
		Children: []model.ChildDictionary{
			{Key: "testCh1", Name: "dziecko 1 updated", Content: content},
			{Key: "testCh3", Name: "dziecko 3 new", Content: content},
		},
	}

	if err := service.UpdateParent(&p); err != nil {
		t.Errorf("Problem with save parent %v", err)
	}
}

func TestDictionaryService_LoadTranslated(t *testing.T) {

	res, err := service.LoadTranslated("uw", "AbsenceType", "", language.MustParse("en"))
	assert.NoError(t, err)

	fmt.Println(res)

}
