package gonion

import (
	"net/http"

	"golang.org/x/net/context"
)

type Adapter struct {
	Ctx     context.Context
	Handler Handler
}

func (ca *Adapter) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ca.Handler.ServeHTTPContext(ca.Ctx, rw, req)
}
