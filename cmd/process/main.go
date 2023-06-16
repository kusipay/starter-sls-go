package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/kusipay/starter-sls-go/pkg/middleware"
	"github.com/mefellows/vesper"
)

// Response
type Response events.SQSEventResponse

// Body
type Body struct {
	Id    int  `json:"id"`
	Error bool `json:"error"`
}

// Character
type Character struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, event events.SQSEvent) (Response, error) {
	region := os.Getenv("ENV_REGION")
	tableName := os.Getenv("TABLE")

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return Response{}, err
	}

	client := dynamodb.NewFromConfig(cfg)

	for _, record := range event.Records {
		fmt.Println(record)

		body := new(Body)
		_ = json.Unmarshal([]byte(record.Body), body)

		if body.Error {
			return Response{}, errors.New("fatal error")
		}

		character, err := getCharacter(body.Id)
		if err != nil {
			return Response{}, err
		}

		item, err := attributevalue.MarshalMap(map[string]interface{}{
			"id":     strconv.Itoa(character.Id),
			"name":   character.Name,
			"status": character.Status,
		})
		if err != nil {
			return Response{}, err
		}

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: &tableName,
			Item:      item,
		})

	}

	return Response{}, err
}

func getCharacter(id int) (*Character, error) {
	url := fmt.Sprintf("https://rickandmortyapi.com/api/character/%d", id)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("not found")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	character := new(Character)
	json.Unmarshal(body, character)

	return character, nil
}

func main() {
	v := vesper.New(Handler).
		Use(middleware.LogMiddleware())

	v.Start()
}
