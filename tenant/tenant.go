package tenant

import (
	"context"
	"net/http"
	"strconv"
	"strings"
)

var (
	headerName      = "X-Tenant-ID"
	mergeHeaderName = "X-Tenant-Merge-Default"
)

type key int

const tenantKey key = 0

type Tenant struct {
	Name         string
	MergeDefault bool
}

func HeaderName(tenantHeader, mergeHeader string) {
	headerName = tenantHeader
	mergeHeaderName = mergeHeader
}

func FromRequest(req *http.Request) Tenant {

	t := req.Header.Get(headerName)
	if strings.ToLower(t) == "default" {
		t = "" // compatibility
	}
	merge, _ := strconv.ParseBool(req.Header.Get(mergeHeaderName))

	return Tenant{
		Name:         t,
		MergeDefault: merge,
	}
}

func NewContext(ctx context.Context, tenant Tenant) context.Context {
	return context.WithValue(ctx, tenantKey, tenant)
}

func FromContext(ctx context.Context) (Tenant, bool) {
	t, ok := ctx.Value(tenantKey).(Tenant)
	return t, ok
}
