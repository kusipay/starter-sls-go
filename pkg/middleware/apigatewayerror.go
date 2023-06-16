package middleware

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mefellows/vesper"
)

func handleError(err error) events.APIGatewayProxyResponse {
	message := err.Error()

	bytes, _ := json.Marshal(map[string]string{
		"message": "unexpected error: " + message,
	})

	return events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body:       string(bytes),
	}
}

// ApiGatewayErrorMiddleware middleware to manage unhandled errors.
func ApiGatewayErrorMiddleware() vesper.Middleware {
	return func(next vesper.LambdaFunc) vesper.LambdaFunc {
		return func(ctx context.Context, event any) (any, error) {
			response, err := next(ctx, event)
			if err == nil {
				return response, err
			}

			return handleError(err), nil
		}
	}
}
