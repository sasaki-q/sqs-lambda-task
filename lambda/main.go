package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(_ context.Context, e events.SQSEvent) error {
	for _, message := range e.Records {
		fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)
	}
	return nil
}

func main() {
	lambda.Start(Handler)
}
