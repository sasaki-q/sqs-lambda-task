package main

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	lambda "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	sqs "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type InfraStackProps struct {
	awscdk.StackProps
}

func NewInfraStack(scope constructs.Construct, id string, props *InfraStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)
	q := sqs.NewQueue(stack, jsii.String("queue"), &sqs.QueueProps{
		QueueName: jsii.String("queue"),
	})
	ar := awsiam.NewRole(stack, jsii.String("assume-role"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("lambda.amazonaws.com"), nil),
		RoleName:  jsii.String("lambda-assume-role"),
	})
	ar.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   &[]*string{jsii.String("sqs:*"), jsii.String("logs:*")},
		Effect:    awsiam.Effect_ALLOW,
		Resources: &[]*string{jsii.String("*")},
	}))
	l := lambda.NewFunction(stack, jsii.String("lambda"), &lambda.FunctionProps{
		FunctionName: jsii.String("function"),
		MemorySize:   jsii.Number(512),
		Handler:      jsii.String("handler"),
		Runtime:      lambda.Runtime_GO_1_X(),
		Code:         lambda.Code_FromAsset(jsii.String("./bin/"), nil),
		Role:         ar,
	})
	l.AddEventSourceMapping(jsii.String("es-mapping"), &lambda.EventSourceMappingOptions{
		EventSourceArn: q.QueueArn(),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)
	env := app.Node().TryGetContext(jsii.String("ENV"))
	if env == nil {
		panic("please pass context")
	}

	awscdk.Tags_Of(app).Add(jsii.String("Project"), jsii.String("cdk"), nil)
	awscdk.Tags_Of(app).Add(jsii.String("Env"), jsii.String(fmt.Sprintf("%s", env)), nil)

	NewInfraStack(app, "InfraStack", &InfraStackProps{
		awscdk.StackProps{
			Synthesizer: awscdk.NewDefaultStackSynthesizer(
				&awscdk.DefaultStackSynthesizerProps{
					FileAssetsBucketName: jsii.String("sasaki-2024-01-10"),
				},
			),
		},
	})

	app.Synth(nil)
}
