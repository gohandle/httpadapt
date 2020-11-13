package httpadapt

import (
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestResponse(t *testing.T) {
	for i, c := range []struct {
		rec    *httptest.ResponseRecorder
		expErr error
		exp    events.APIGatewayProxyResponse
	}{
		{rec: buildResp(t), exp: events.APIGatewayProxyResponse{
			StatusCode:        200,
			MultiValueHeaders: map[string][]string{"Content-Type": {"text/plain; charset=utf-8"}},
		}},

		{rec: buildResp(t, respWriteHeader(400)), exp: events.APIGatewayProxyResponse{
			StatusCode:        400,
			MultiValueHeaders: map[string][]string{"Content-Type": {"text/plain; charset=utf-8"}},
		}},

		{rec: buildResp(t, respWrite([]byte{0x75})), exp: events.APIGatewayProxyResponse{
			StatusCode:        200,
			Body:              "u",
			MultiValueHeaders: map[string][]string{"Content-Type": {"text/plain; charset=utf-8"}},
			IsBase64Encoded:   false,
		}},
		{rec: buildResp(t, respWrite([]byte("\xe2\x28\xa1"))), exp: events.APIGatewayProxyResponse{
			StatusCode:        200,
			Body:              "4iih",
			MultiValueHeaders: map[string][]string{"Content-Type": {"text/plain; charset=utf-8"}},
			IsBase64Encoded:   true,
		}},

		{rec: buildResp(t, //write header after write body
			respWrite([]byte("foo")),
			respHeaderSet(map[string][]string{"foo": {"bar"}}),
		), exp: events.APIGatewayProxyResponse{
			StatusCode:        200,
			Body:              "foo",
			MultiValueHeaders: map[string][]string{"Content-Type": {"text/plain; charset=utf-8"}},
		}},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			out := ProxyResponse(c.rec)
			if !reflect.DeepEqual(out, c.exp) {
				t.Fatalf("got: %v exp: %v", out, c.exp)
			}
		})
	}
}

type respTrans func(testing.TB, *httptest.ResponseRecorder)

func buildResp(tb testing.TB, trs ...respTrans) (w *httptest.ResponseRecorder) {
	w = httptest.NewRecorder()
	for _, tr := range trs {
		tr(tb, w)
	}
	return
}

func respWrite(b []byte) respTrans {
	return func(tb testing.TB, w *httptest.ResponseRecorder) {
		if _, err := w.Write(b); err != nil {
			tb.Fatal(err)
		}
	}
}

func respWriteHeader(c int) respTrans {
	return func(tb testing.TB, w *httptest.ResponseRecorder) {
		w.WriteHeader(c)
	}
}

func respHeaderSet(m map[string][]string) respTrans {
	return func(tb testing.TB, w *httptest.ResponseRecorder) {
		for k, vs := range m {
			for _, v := range vs {
				w.Header().Add(k, v)
			}
		}
	}
}
