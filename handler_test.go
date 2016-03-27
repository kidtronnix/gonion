package gonion

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/net/context"

	"github.com/stretchr/testify/assert"
)

func Middleware1(c Handler) Handler {
	return HandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
		ctx = context.WithValue(ctx, "args", []string{"hello"})
		c.ServeHTTPContext(ctx, rw, req)
	})
}

func Middleware2(c Handler) Handler {
	return HandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
		args := ctx.Value("args").([]string)
		ctx = context.WithValue(ctx, "args", append(args, "world"))
		c.ServeHTTPContext(ctx, rw, req)
	})
}

func Response(c Handler) Handler {
	return HandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
		args := ctx.Value("args").([]string)
		str := strings.Join(args, " ")
		fmt.Fprint(rw, str)
	})
}

func TestGonionHandler(t *testing.T) {
	assert := assert.New(t)

	handler := NewHTTPHandler(
		Middleware1,
		Middleware2,
		Response,
	)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "http://example.com", nil)
	assert.NoError(err)

	handler.ServeHTTP(w, r)
	assert.Equal("hello world", w.Body.String())
}
