package configuration

import (
	"fmt"
	"testing"
)

func TestConfigurationRepository_LoadAllShort(t *testing.T) {

	short, err := confRepository.LoadAllShort("")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", short)

}
