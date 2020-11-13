package httpadapt

import (
	"context"
	"net/http"
	"testing"
)

func TestAdapt(t *testing.T) {
	ctx := context.Background()
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context() != ctx {
			t.Fatalf("should have been passed the context")
		}

		w.WriteHeader(415)
	})

	out, err := New(h).ProxyWithContext(ctx, event("/hello", "GET"))
	if err != nil {
		t.Fatalf("got: %v", err)
	}

	if out.StatusCode != 415 {
		t.Fatalf("got: %v", out.StatusCode)
	}
}
