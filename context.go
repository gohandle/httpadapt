package httpadapt

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

// WithContext returns a context with gateway context and lamda context added (if any)
func WithContext(ctx context.Context, apiGwRequest events.APIGatewayProxyRequest) context.Context {
	lc, _ := lambdacontext.FromContext(ctx)
	return context.WithValue(ctx, contextKey{}, requestContext{
		lambdaContext:       lc,
		gatewayProxyContext: apiGwRequest.RequestContext,
		stageVars:           apiGwRequest.StageVariables,
	})
}

// withContext returns the httpRequest with gateway context and lambda context added
func withContext(ctx context.Context, req *http.Request, apiGwRequest events.APIGatewayProxyRequest) *http.Request {
	return req.WithContext(WithContext(ctx, apiGwRequest))
}

// GetAPIGatewayContextFromContext retrieve APIGatewayProxyRequestContext from context.Context
func GetAPIGatewayContextFromContext(ctx context.Context) (events.APIGatewayProxyRequestContext, bool) {
	v, ok := ctx.Value(contextKey{}).(requestContext)
	return v.gatewayProxyContext, ok
}

// GetRuntimeContextFromContext retrieve Lambda Runtime Context from context.Context
func GetRuntimeContextFromContext(ctx context.Context) (*lambdacontext.LambdaContext, bool) {
	v, ok := ctx.Value(contextKey{}).(requestContext)
	return v.lambdaContext, ok
}

// GetStageVarsFromContext retrieve stage variables from context
func GetStageVarsFromContext(ctx context.Context) (map[string]string, bool) {
	v, ok := ctx.Value(contextKey{}).(requestContext)
	return v.stageVars, ok
}

type contextKey struct{}

type requestContext struct {
	lambdaContext       *lambdacontext.LambdaContext
	gatewayProxyContext events.APIGatewayProxyRequestContext
	stageVars           map[string]string
}
