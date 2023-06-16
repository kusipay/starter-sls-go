package main

import (
	"context"
	"encoding/json"

	"github.com/kusipay/starter-sls-go/pkg/middleware"
	"github.com/mefellows/vesper"
)

type Person struct {
	Name string
	Age  int
}

func Handler(ctx context.Context, event interface{}) (map[string]string, error) {
	bytes, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	result := new(map[string]string)

	err = json.Unmarshal(bytes, result)
	if err != nil {
		return nil, err
	}

	return *result, nil
}

func main() {
	v := vesper.New(Handler).Use(middleware.LogMiddleware())

	v.Start()
}
