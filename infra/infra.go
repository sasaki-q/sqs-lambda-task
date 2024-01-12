package main

import (
	"fmt"
	"infra/resources"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	ecr "github.com/aws/aws-cdk-go/awscdk/v2/awsecr"
	ecs "github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	lambda "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	logs "github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
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

	r := resources.NewVpc(stack)
	repository := ecr.Repository_FromRepositoryName(stack, jsii.String("ecr"), jsii.String("task"))
	logBucket := awss3.Bucket_FromBucketName(stack, jsii.String("bucket"), jsii.String("sasaki-2024-01-10"))
	logGroup := logs.LogGroup_FromLogGroupName(stack, jsii.String("log"), jsii.String("log-group"))

	cluster := ecs.NewCluster(stack, jsii.String("cluster"), &ecs.ClusterProps{
		ClusterName:                    jsii.String("cluster"),
		EnableFargateCapacityProviders: jsii.Bool(true),
		Vpc:                            r.Vpc,
		ExecuteCommandConfiguration: &ecs.ExecuteCommandConfiguration{
			LogConfiguration: &ecs.ExecuteCommandLogConfiguration{
				CloudWatchLogGroup:          logGroup,
				CloudWatchEncryptionEnabled: jsii.Bool(true),
				S3Bucket:                    logBucket,
				S3EncryptionEnabled:         jsii.Bool(true),
				S3KeyPrefix:                 jsii.String("log"),
			},
			Logging: ecs.ExecuteCommandLogging_OVERRIDE,
		},
	})

	task := ecs.NewTaskDefinition(stack, jsii.String("task"), &ecs.TaskDefinitionProps{
		Compatibility:   ecs.Compatibility_FARGATE,
		Cpu:             jsii.String("256"),
		MemoryMiB:       jsii.String("512"),
		Family:          jsii.String("family"),
		RuntimePlatform: &ecs.RuntimePlatform{CpuArchitecture: ecs.CpuArchitecture_X86_64()},
	})

	task.AddContainer(jsii.String("container"),
		&ecs.ContainerDefinitionOptions{
			Image:          ecs.ContainerImage_FromEcrRepository(repository, jsii.String("v0.1")),
			Cpu:            jsii.Number(256),
			MemoryLimitMiB: jsii.Number(512),
			ContainerName:  jsii.String("container"),
			Logging:        ecs.LogDriver_AwsLogs(&ecs.AwsLogDriverProps{StreamPrefix: jsii.String("cdklog"), LogGroup: logGroup}),
			Environment:    &map[string]*string{},
		},
	)

	dq := sqs.NewQueue(stack, jsii.String("dead-letter-queue"), &sqs.QueueProps{QueueName: jsii.String("dead-letter-queue")})
	q := sqs.NewQueue(stack, jsii.String("queue"), &sqs.QueueProps{
		QueueName:       jsii.String("queue"),
		DeadLetterQueue: &sqs.DeadLetterQueue{MaxReceiveCount: jsii.Number(1), Queue: dq},
	})
	ar := awsiam.NewRole(stack, jsii.String("assume-role"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("lambda.amazonaws.com"), nil),
		RoleName:  jsii.String("lambda-assume-role"),
	})
	ar.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   &[]*string{jsii.String("sqs:*"), jsii.String("logs:*"), jsii.String("ecs:RunTask"), jsii.String("iam:PassRole")},
		Effect:    awsiam.Effect_ALLOW,
		Resources: &[]*string{jsii.String("*")},
	}))
	l := lambda.NewFunction(stack, jsii.String("lambda"), &lambda.FunctionProps{
		FunctionName: jsii.String("function"),
		MemorySize:   jsii.Number(512),
		Handler:      jsii.String("infra/bin/handler"),
		Runtime:      lambda.Runtime_GO_1_X(),
		Code:         lambda.AssetCode_FromAsset(jsii.String("./bin/handler.zip"), nil),
		Role:         ar,
		Environment: &map[string]*string{
			"Cluster":        cluster.ClusterName(),
			"Task":           task.TaskDefinitionArn(),
			"Container":      jsii.String("container"),
			"Subnets":        jsii.String(r.SubnetIds),
			"SecurityGroups": jsii.String(r.SecurityGroupIds),
		},
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
