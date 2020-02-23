package service

import (
	"fmt"
	"testing"
)

func TestChildRepository_LoadChildren(t *testing.T) {
	ch, err := chRepository.LoadChildren("uw", "AbsenceType", "")

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%v\n", ch)

}
