package httpadapt

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/aws/aws-lambda-go/events"
)

// Adapter adapts http.Handler implementations.
type Adapter struct {
	h http.Handler

	stripBasePath       string
	customServerAddress string
}

// New inits an adapter
func New(h http.Handler, opts ...Option) (a *Adapter) {
	a = &Adapter{h: h}
	for _, opt := range opts {
		opt(a)
	}

	return
}

// ProxyWithContext receives context and an API Gateway proxy event,
// transforms them into an http.Request object, and sends it to the http.Handler for routing.
// It returns a proxy response object generated from the http.ResponseWriter.
func (a *Adapter) ProxyWithContext(
	ctx context.Context,
	ev events.APIGatewayProxyRequest,
) (out events.APIGatewayProxyResponse, err error) {
	req, err := a.EventToRequest(ev)
	if err != nil {
		return out, fmt.Errorf("failed to create request from event: %w", err)
	}

	rec := httptest.NewRecorder()
	a.h.ServeHTTP(rec, withContext(ctx, req, ev)) // call the implemention
	return ProxyResponse(rec), nil
}
