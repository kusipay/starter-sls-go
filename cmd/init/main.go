package main

import (
	"bytes"
	"context"
	"encoding/json"
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/kusipay/starter-sls-go/pkg/middleware"
	"github.com/mefellows/vesper"
)

// Body
type Body struct {
	Error bool `json:"error"`
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	payload := new(Body)
	_ = json.Unmarshal([]byte(event.Body), payload)

	err := PostToSQS(payload.Error)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 418}, nil
	}

	var buf bytes.Buffer

	body, err := json.Marshal(map[string]interface{}{
		"message": "Say hello to my little friends",
	})

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 404}, nil
	}

	json.HTMLEscape(&buf, body)

	resp := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       buf.String(),
	}

	return resp, nil
}

func randomizer(n int) int {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	return r.Intn(n) + 1
}

// PostToSQS post message to sqs
func PostToSQS(withError bool) error {
	region := os.Getenv("ENV_REGION")
	queueUrl := os.Getenv("QUEUE")

	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return err
	}

	client := sqs.NewFromConfig(cfg)

	id := randomizer(826)

	message, _ := json.Marshal(map[string]any{
		"id":    id,
		"error": withError,
	})

	_, err = client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &queueUrl,
		MessageBody: aws.String(string(message)),
	})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	v := vesper.New(Handler).
		Use(middleware.LogMiddleware()).
		Use(middleware.HeadersMiddleware(middleware.Cors, middleware.Json)).
		Use(middleware.ApiGatewayErrorMiddleware())

	v.Start()
}
