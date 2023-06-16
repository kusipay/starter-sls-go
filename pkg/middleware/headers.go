package middleware

import (
	"context"
	"reflect"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mefellows/vesper"
)

const (
	// Cors sets cors headers.
	Cors = "cors"

	// Json sets application/json.
	Json = "json"
)

func addHeaders(response events.APIGatewayProxyResponse, headers ...string) events.APIGatewayProxyResponse {
	if response.Headers == nil {
		response.Headers = map[string]string{}
	}

	for _, header := range headers {
		if header == Cors {
			response.Headers["Access-Control-Allow-Origin"] = "*"
		}

		if header == Json {
			response.Headers["Content-Type"] = "application/json"
		}
	}

	return response
}

// HeadersMiddleware middleware to add headers in the apigateway response.
func HeadersMiddleware(headers ...string) vesper.Middleware {
	return func(next vesper.LambdaFunc) vesper.LambdaFunc {
		return func(ctx context.Context, event any) (any, error) {
			response, err := next(ctx, event)

			var newresp any

			t := reflect.TypeOf(response)
			if t.Kind() == reflect.Pointer {
				t = t.Elem()
			}

			if t.Kind() == reflect.Struct && t.ConvertibleTo(reflect.TypeOf(events.APIGatewayProxyResponse{})) {
				resp := reflect.ValueOf(response).Convert(reflect.TypeOf(events.APIGatewayProxyResponse{})).
					Interface().(events.APIGatewayProxyResponse)

				temp := addHeaders(resp)

				if t.Kind() == reflect.Pointer {
					newresp = &temp
				} else {
					newresp = temp
				}
			} else {
				newresp = response
			}

			return newresp, err
		}
	}
}
