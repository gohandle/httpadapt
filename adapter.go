package httpadapt

import "net/http"

// Adapter adapts http.Handler implementations.
type Adapter struct {
	h http.Handler
}

// New inits an adapter
func New(h http.Handler, opts ...Option) (a *Adapter) {
	a = &Adapter{h}
	for _, opt := range opts {
		opt(a)
	}

	return
}