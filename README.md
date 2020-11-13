# httpadapt
Adapt a http.Handler to handle API Gateway Proxy events in AWS Lambda. This library started as 
simple fork of the  blessed `https://github.com/awslabs/aws-lambda-go-api-proxy` library but in the
end so much was changed that only very little of the original codebase remains.

## features
- This library only depends on the standard library and `github.com/aws/aws-lambda-go`
- We removed all support for non-context based HTT handling to vastly simply the code base
- Only supports the standard http.Handler interface to vastly simplify the code
- Query parameters aren now ordered deterministicly instead
- Instead of writing our own ResponseWriter we use the battle tested httptest.ResponseRecorder 
- Well tested with coverage of over 95%

## usage

```Go
package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gohandle/httpadapt"
)

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, %v", r.URL)
}

func main() {
	lambda.Start(
		httpadapt.New(http.HandlerFunc(handle)).ProxyWithContext)
}
```

## backlog
- [x] Add a functional option to configure stripbasepath
- [x] Add a functional option for CustomHostVariable env
- [ ] Test errors, possiblty with a package error type
- [ ] Consider the v2 api format
- [x] Prevent header from being edited after writing with Write, else it will work on lambda but 
      not in a real server