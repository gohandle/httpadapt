package httpadapt

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// EventToRequest converts an API Gateway proxy event into an http.Request object.
// Returns the populated request maintaining headers
func (a *Adapter) EventToRequest(ev events.APIGatewayProxyRequest) (*http.Request, error) {
	decodedBody := []byte(ev.Body)
	if ev.IsBase64Encoded {
		base64Body, err := base64.StdEncoding.DecodeString(ev.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64 body: %w", err)
		}
		decodedBody = base64Body
	}

	path := ev.Path
	if a.stripBasePath != "" {
		if strings.HasPrefix(path, a.stripBasePath) {
			path = strings.Replace(path, a.stripBasePath, "", 1)
		}
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	serverAddress := "https://" + ev.RequestContext.DomainName
	if a.customServerAddress != "" {
		serverAddress = a.customServerAddress
	}
	path = serverAddress + path

	if len(ev.MultiValueQueryStringParameters) > 0 {
		queryString := ""

		keys := make([]string, 0, len(ev.MultiValueQueryStringParameters))
		for q := range ev.MultiValueQueryStringParameters {
			keys = append(keys, q)
		}

		sort.Strings(keys)
		for _, q := range keys {
			for _, v := range ev.MultiValueQueryStringParameters[q] {
				if queryString != "" {
					queryString += "&"
				}
				queryString += url.QueryEscape(q) + "=" + url.QueryEscape(v)
			}
		}
		path += "?" + queryString
	} else if len(ev.QueryStringParameters) > 0 {

		keys := make([]string, 0, len(ev.QueryStringParameters))
		for q := range ev.QueryStringParameters {
			keys = append(keys, q)
		}

		// Support `QueryStringParameters` for backward compatibility.
		// https://github.com/awslabs/aws-lambda-go-api-proxy/issues/37
		queryString := ""
		sort.Strings(keys)
		for _, q := range keys {
			if queryString != "" {
				queryString += "&"
			}
			queryString += url.QueryEscape(q) + "=" + url.QueryEscape(ev.QueryStringParameters[q])
		}
		path += "?" + queryString
	}

	req, err := http.NewRequest(
		strings.ToUpper(ev.HTTPMethod),
		path,
		bytes.NewReader(decodedBody),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if ev.MultiValueHeaders != nil {
		for k, values := range ev.MultiValueHeaders {
			for _, value := range values {
				req.Header.Add(k, value)
			}
		}
	} else {
		for h := range ev.Headers {
			req.Header.Add(h, ev.Headers[h])
		}
	}

	req.RequestURI = req.URL.RequestURI()
	return req, nil
}
