# httpadapt
Adapt a http.Handler to a Lambda Gateway Event handler

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

## features
- Stdlib and lambda event deps only
- Only supports context based handling
- Deterministic query params order
- Battle tested httptest.ResponseRecorder to record the response

## backlog
- [x] Add a functional option to configure stripbasepath
- [x] Add a functional option for CustomHostVariable env
- [ ] Test errors, possiblty with a package error type
- [ ] Consider the v2 api format
- [x] Prevent header from being edited after writing with Write, else it will work on lambda but 
      not in a real server