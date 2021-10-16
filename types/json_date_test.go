package types

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestJsonDate_Format(t *testing.T) {

	now := time.Now()

	marshal, err := json.Marshal(now)
	if err != nil {
		return
	}

	fmt.Println(string(marshal))

}
