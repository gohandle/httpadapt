package httpadapt

import (
	"errors"
	"reflect"
	"strconv"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestResponse(t *testing.T) {
	for i, c := range []struct {
		w      *ResponseWriter
		expErr error
		exp    events.APIGatewayProxyResponse
	}{
		{w: buildResp(t)},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			out, err := c.w.ProxyResponse()
			if !errors.Is(err, c.expErr) {
				t.Fatalf("got: %v exp: %v", err, c.expErr)
			}

			if !reflect.DeepEqual(out, c.exp) {
				t.Fatalf("got: %v exp: %v", out, c.exp)
			}
		})
	}
}

type respTrans func(testing.TB, *ResponseWriter)

func buildResp(tb testing.TB, trs ...respTrans) (w *ResponseWriter) {
	w = NewResponseWriter()
	for _, tr := range trs {
		tr(tb, w)
	}
	return
}

func respWrite(b []byte) respTrans {
	return func(tb testing.TB, w *ResponseWriter) {
		if _, err := w.Write(b); err != nil {
			tb.Fatal(err)
		}
	}
}

func respWriteHeader(c int) respTrans {
	return func(tb testing.TB, w *ResponseWriter) {
		w.WriteHeader(c)
	}
}

func respHeaderSet(m map[string][]string) respTrans {
	return func(tb testing.TB, w *ResponseWriter) {
		for k, vs := range m {
			for _, v := range vs {
				w.Header().Add(k, v)
			}
		}
	}
}
