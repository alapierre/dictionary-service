package calendar

import (
	"context"
	"dictionaries-service/tenant"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_service_Load(t *testing.T) {

	var from = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	var to = time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC)

	ctx := tenant.NewContext(context.Background(), tenant.Tenant{})

	res, err := calendarService.LoadByTypeAndRange(ctx, "holiday", from, to)

	assert.NoError(t, err)
	fmt.Printf("%#v\n", res)
}
