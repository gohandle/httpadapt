package httpadapt

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestEventToRequest(t *testing.T) {
	for i, c := range []struct {
		ev        events.APIGatewayProxyRequest
		adapter   *Adapter
		expErr    error
		expMethod string
		expURL    *url.URL
		expBody   []byte
		expHeader http.Header
		expURI    string
	}{
		{ // zero value
			adapter: New(nil),
			expURI:  "/", expURL: &url.URL{Scheme: "https", Path: "/"},
			expHeader: http.Header{}, expMethod: "GET",
		},
		{ // with lower case post
			adapter: New(nil),
			ev:      event("/hello", "post"),
			expURI:  "/hello", expURL: &url.URL{Scheme: "https", Path: "/hello"},
			expHeader: http.Header{}, expMethod: "POST",
		},
		{ // with binary body
			adapter: New(nil),
			ev: event("/hello", "PUT",
				body([]byte{0x01}),
			),
			expURI: "/hello", expURL: &url.URL{Scheme: "https", Path: "/hello"},
			expHeader: http.Header{}, expMethod: "PUT", expBody: []byte{0x01},
		},
		{ // both types of query strings
			adapter: New(nil),
			ev: event("/hello", "DElEtE",
				multiValueQueryStringParameters(map[string][]string{
					"hello": {"1"},
					"world": {"2", "3"},
				}),
				queryStringParameters(map[string]string{
					"hello": "1",
					"world": "2",
				}),
			),
			expURI: "/hello?hello=1&world=2&world=3", expHeader: http.Header{}, expMethod: "DELETE",
			expURL: &url.URL{Scheme: "https", Path: "/hello", RawQuery: "hello=1&world=2&world=3"},
		},

		{ // both types of query strings
			adapter: New(nil),
			ev: event("/hello", "GET",
				queryStringParameters(map[string]string{
					"hello": "1",
					"world": "2",
				}),
			),
			expURI: "/hello?hello=1&world=2", expHeader: http.Header{}, expMethod: "GET",
			expURL: &url.URL{Scheme: "https", Path: "/hello", RawQuery: "hello=1&world=2"},
		},

		{ // both types of query strings
			adapter: New(nil),
			ev: event("/hello", "GET",
				multiValueQueryStringParameters(map[string][]string{
					"hello": {"1"},
					"world": {"2", "3"},
				}),
			),
			expURI: "/hello?hello=1&world=2&world=3", expHeader: http.Header{}, expMethod: "GET",
			expURL: &url.URL{Scheme: "https", Path: "/hello", RawQuery: "hello=1&world=2&world=3"},
		},

		{ // both types of query strings
			adapter: New(nil),
			ev: event("/hello", "GET",
				headers(map[string]string{
					"hello": "1",
					"world": "2",
				}),
			),
			expURI: "/hello", expHeader: http.Header{"Hello": {"1"}, "World": {"2"}}, expMethod: "GET",
			expURL: &url.URL{Scheme: "https", Path: "/hello"},
		},

		{ // single headers headers
			adapter: New(nil),
			ev: event("/hello", "GET",
				headers(map[string]string{
					"hello": "1",
					"world": "2",
				}),
			),
			expURI: "/hello", expHeader: http.Header{"Hello": {"1"}, "World": {"2"}}, expMethod: "GET",
			expURL: &url.URL{Scheme: "https", Path: "/hello"},
		},

		{ // multi-value headers
			adapter: New(nil),
			ev: event("/hello", "GET",
				multiValueHeaders(map[string][]string{
					"hello": {"1"},
					"world": {"2", "3"},
				}),
			),
			expURI: "/hello", expHeader: http.Header{"Hello": {"1"}, "World": {"2", "3"}}, expMethod: "GET",
			expURL: &url.URL{Scheme: "https", Path: "/hello"},
		},

		{ // strip base path
			adapter: New(nil, StripBasePath("/base")),
			ev:      event("/base/hello", "GET"),
			expURI:  "/hello", expHeader: http.Header{}, expMethod: "GET",
			expURL: &url.URL{Scheme: "https", Path: "/hello"},
		},

		{ // custom host path
			adapter: New(nil, CustomHost("http://custom.host.com")),
			ev:      event("/hello", "GET"),
			expURI:  "/hello", expHeader: http.Header{}, expMethod: "GET",
			expURL: &url.URL{Scheme: "http", Host: "custom.host.com", Path: "/hello"},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			req, err := c.adapter.EventToRequest(c.ev)
			if !errors.Is(err, c.expErr) {
				t.Fatalf("got: %v exp: %v", err, c.expErr)
			}

			if act := req.Method; act != c.expMethod {
				t.Fatalf("got: %v exp: %v", act, c.expMethod)
			}

			if act := req.URL; !reflect.DeepEqual(act, c.expURL) {
				t.Fatalf("got: %v exp: %v", act, c.expURL)
			}

			if act, err := ioutil.ReadAll(req.Body); !bytes.Equal(act, c.expBody) {
				t.Fatalf("got: %v exp: %v (%v)", act, c.expBody, err)
			}

			if act := req.Header; !reflect.DeepEqual(act, c.expHeader) {
				t.Fatalf("got: %v exp: %v", act, c.expHeader)
			}

			if act := req.RequestURI; act != c.expURI {
				t.Fatalf("got: %v exp: %v", act, c.expURI)
			}
		})
	}
}

// opt configures test events
type opt func(ev events.APIGatewayProxyRequest) events.APIGatewayProxyRequest

// event creates an empty test event
func event(path string, method string, opts ...opt) events.APIGatewayProxyRequest {
	ev := events.APIGatewayProxyRequest{
		Path:       path,
		HTTPMethod: method,
	}

	for _, o := range opts {
		ev = o(ev)
	}

	return ev
}

func body(b []byte) opt {
	return func(ev events.APIGatewayProxyRequest) events.APIGatewayProxyRequest {
		ev.Body = base64.StdEncoding.EncodeToString(b)
		ev.IsBase64Encoded = true
		return ev
	}
}

func multiValueQueryStringParameters(m map[string][]string) opt {
	return func(ev events.APIGatewayProxyRequest) events.APIGatewayProxyRequest {
		ev.MultiValueQueryStringParameters = m
		return ev
	}
}

func queryStringParameters(m map[string]string) opt {
	return func(ev events.APIGatewayProxyRequest) events.APIGatewayProxyRequest {
		ev.QueryStringParameters = m
		return ev
	}
}

func headers(m map[string]string) opt {
	return func(ev events.APIGatewayProxyRequest) events.APIGatewayProxyRequest {
		ev.Headers = m
		return ev
	}
}

func multiValueHeaders(m map[string][]string) opt {
	return func(ev events.APIGatewayProxyRequest) events.APIGatewayProxyRequest {
		ev.MultiValueHeaders = m
		return ev
	}
}
