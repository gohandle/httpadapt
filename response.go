package httpadapt

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"unicode/utf8"

	"github.com/aws/aws-lambda-go/events"
)

// ProxyResponse returns the resulting response to handling the API Gateway Proxy event
func ProxyResponse(rec *httptest.ResponseRecorder) (out events.APIGatewayProxyResponse) {

	// if the content type header is not set when we write the body we try to
	// detect one and set it by default. If the content type cannot be detected
	// it is automatically set to "application/octet-stream" by the
	// DetectContentType method
	if rec.Result().Header.Get("Content-Type") == "" {
		rec.Result().Header.Add("Content-Type", http.DetectContentType(
			rec.Body.Bytes(),
		))
	}

	out.StatusCode = rec.Result().StatusCode
	out.MultiValueHeaders = rec.Result().Header
	out.Body = string(rec.Body.Bytes())
	if !utf8.Valid(rec.Body.Bytes()) {
		out.Body = base64.StdEncoding.EncodeToString(rec.Body.Bytes())
		out.IsBase64Encoded = true
	}

	return
}
