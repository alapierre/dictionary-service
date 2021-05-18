package tenant

import (
	"context"
	"net/http"
	"strings"
)

var headerName = "X-Tenant-ID"

type key int

const tenantKey key = 0

func HeaderName(name string) {
	headerName = name
}

func FromRequest(req *http.Request) string {
	t := req.Header.Get(headerName)
	if strings.ToLower(t) == "default" {
		t = "" // compatibility
	}
	return t
}

func NewContext(ctx context.Context, tenant string) context.Context {
	return context.WithValue(ctx, tenantKey, tenant)
}

func FromContext(ctx context.Context) (string, bool) {
	userIP, ok := ctx.Value(tenantKey).(string)
	return userIP, ok
}
