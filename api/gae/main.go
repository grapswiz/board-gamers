package api

import (
	"net/http"

	"github.com/favclip/ucon"
	"github.com/favclip/ucon/swagger"
	"google.golang.org/appengine"
)

func init() {
	ucon.Middleware(UseAppengineContext)
	ucon.Orthodox()
	ucon.Middleware(swagger.RequestValidator())

	ucon.DefaultMux.Prepare()
	http.Handle("/", ucon.DefaultMux)
}

func UseAppengineContext(b *ucon.Bubble) error {
	if b.Context == nil {
		b.Context = appengine.NewContext(b.R)
	} else {
		b.Context = appengine.WithContext(b.Context, b.R)
	}
	b.R = b.R.WithContext(b.Context)

	return b.Next()
}
