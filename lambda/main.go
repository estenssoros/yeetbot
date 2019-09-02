package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func handleLambda(ctx context.Context, req map[string]interface{}) error {
	fmt.Println(req)
	return nil
}

func main() {
	lambda.Start(handleLambda)
}
