package httpadapt

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// EventToRequest converts an API Gateway proxy event into an http.Request object.
// Returns the populated request maintaining headers
func (a *Adapter) EventToRequest(ev events.APIGatewayProxyRequest) (req *http.Request, err error) {
	body := []byte(ev.Body)
	if ev.IsBase64Encoded {
		body, err = base64.StdEncoding.DecodeString(ev.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64 body: %w", err)
		}
	}

	p := ev.Path
	if a.stripBasePath != "" && strings.HasPrefix(p, a.stripBasePath) {
		p = strings.Replace(p, a.stripBasePath, "", 1)
	}

	if !path.IsAbs(p) {
		p = path.Join("/", p)
	}

	serverAddress := "https://" + ev.RequestContext.DomainName
	if a.customServerAddress != "" {
		serverAddress = a.customServerAddress
	}
	p = serverAddress + p

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
		p += "?" + queryString
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
		p += "?" + queryString
	}

	if req, err = http.NewRequest(
		strings.ToUpper(ev.HTTPMethod), p,
		bytes.NewReader(body),
	); err != nil {
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
