package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func Handler(ctx context.Context, e events.SQSEvent) error {
	var count int32 = 1
	cluster := os.Getenv("Cluster")
	taskArn := os.Getenv("Task")
	container := os.Getenv("Container")
	subnets := os.Getenv("Subnets")
	securityGroups := os.Getenv("SecurityGroups")

	if cluster == "" || taskArn == "" || container == "" || subnets == "" || securityGroups == "" {
		log.Fatal("ERROR: environment vairables are missing")
		return nil
	}

	subnetIds := removeDuplicateValue(strings.Split(subnets, ","))
	securityGroupIds := removeDuplicateValue(strings.Split(securityGroups, ","))

	client, err := newClient(ctx)
	if err != nil {
		log.Fatalf("ERROR: cannot create client === %s", err)
		return err
	}

	for _, message := range e.Records {
		op, err := client.RunTask(
			ctx,
			&ecs.RunTaskInput{
				Cluster:        aws.String(cluster),
				TaskDefinition: aws.String(taskArn),
				LaunchType:     types.LaunchTypeFargate,
				NetworkConfiguration: &types.NetworkConfiguration{
					AwsvpcConfiguration: &types.AwsVpcConfiguration{
						AssignPublicIp: types.AssignPublicIpDisabled,
						SecurityGroups: securityGroupIds,
						Subnets:        subnetIds,
					},
				},
				Count: &count,
				Overrides: &types.TaskOverride{
					ContainerOverrides: []types.ContainerOverride{
						{
							Name: aws.String(container),
							Environment: []types.KeyValuePair{
								{Name: aws.String("MessageId"), Value: aws.String(message.MessageId)},
							},
						},
					},
				},
			},
		)

		if err != nil {
			log.Fatalf("ERROR: run task error === %s \n", err)
			return err
		}

		fmt.Print("DEBUG: run task output === ", op, "\n")
	}

	return nil
}

func newClient(ctx context.Context) (*ecs.Client, error) {
	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	return ecs.NewFromConfig(awsConfig), nil
}

func removeDuplicateValue(e []string) []string {
	tmp := make(map[string]bool)
	uniqueIds := []string{}

	for _, id := range e {
		if !tmp[id] {
			tmp[id] = true
			uniqueIds = append(uniqueIds, id)
		}
	}
	return uniqueIds
}

func main() {
	lambda.Start(Handler)
}
