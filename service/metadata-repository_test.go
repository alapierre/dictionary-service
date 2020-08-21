package service

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_dictionaryMetadataRepository_Delete(t *testing.T) {
}

func Test_dictionaryMetadataRepository_Save(t *testing.T) {

}

func Test_dictionaryMetadataRepository_Update(t *testing.T) {

}

func TestDictionaryMetadataRepository_Load(t *testing.T) {

	res, err := metadataRepository.Load("AbsenceType", "")
	assert.NoError(t, err)

	fmt.Printf("%v", res)

}
