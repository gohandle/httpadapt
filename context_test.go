package httpadapt

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

func TestWithContext(t *testing.T) {
	ctx := context.Background()
	ctx = lambdacontext.NewContext(ctx, &lambdacontext.LambdaContext{
		AwsRequestID: "foo",
	})

	req, _ := http.NewRequest("GET", "/", nil)
	req = withContext(ctx, req, events.APIGatewayProxyRequest{
		RequestContext: events.APIGatewayProxyRequestContext{
			AccountID: "foo",
		},
		StageVariables: map[string]string{
			"foo": "bar",
		},
	})

	gwctx, ok := GetAPIGatewayContextFromContext(req.Context())
	if !ok || gwctx.AccountID != "foo" {
		t.Fatalf("got: %v", gwctx)
	}

	svars, ok := GetStageVarsFromContext(req.Context())
	if !ok || svars["foo"] != "bar" {
		t.Fatalf("got: %v", svars)
	}

	lc, ok := GetRuntimeContextFromContext(req.Context())
	if !ok || lc.AwsRequestID != "foo" {
		t.Fatalf("got: %v", lc)
	}
}
