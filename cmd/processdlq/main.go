package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/kusipay/starter-sls-go/pkg/middleware"
	"github.com/kusipay/starter-sls-go/pkg/util"
	"github.com/mefellows/vesper"
)

// Response
type Response events.SQSEventResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, event events.SQSEvent) (Response, error) {
	for _, record := range event.Records {
		util.LogJson("record |", record)
	}

	return Response{}, nil
}

func main() {
	v := vesper.New(Handler).
		Use(middleware.LogMiddleware())

	v.Start()
}
