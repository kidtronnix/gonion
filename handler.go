package gonion

import (
	"net/http"

	"golang.org/x/net/context"
)

type Handler interface {
	ServeHTTPContext(context.Context, http.ResponseWriter, *http.Request)
}

type HandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

func (m HandlerFunc) ServeHTTPContext(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	m(ctx, rw, req)
}

func NewHTTPHandler(ds ...Decorator) http.Handler {
	handler := HandlerFunc(func(context.Context, http.ResponseWriter, *http.Request) {})
	return &Adapter{
		Ctx:     context.Background(),
		Handler: Decorate(handler, ds...),
	}
}
