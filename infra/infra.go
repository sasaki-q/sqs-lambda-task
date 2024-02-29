package main

import (
	"fmt"
	"infra/resources"
	"strings"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type InfraStackProps struct {
	awscdk.StackProps
}

type Props struct {
	Project string
}

const (
	Cidr string = "192.168.0.0/16"

	LogBucketName string = "log-sandbox-sasaki"
	LogGroupName  string = "log-group"
)

func name(a string, e string) string {
	return fmt.Sprintf("%s-%s", strings.ToLower(a), e)
}

func NewInfraStack(scope constructs.Construct, id string, props *InfraStackProps, e Props) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)
	var i resources.IResourceService = &resources.ResourceService{S: stack}

	var (
		vpcName              = name(e.Project, "vpc")
		securityGroupName    = name(e.Project, "sg")
		repositoryName       = name(e.Project, "task-repository")
		clusterName          = name(e.Project, "cluster")
		taskName             = name(e.Project, "task")
		containerName        = name(e.Project, "container")
		queueName            = name(e.Project, "queue")
		deadLetterQueueName  = name(e.Project, "dead-letter-queue")
		lambdaAssumeRoleName = name(e.Project, "lambda-assume-role")
		lambdaPrincipal      = "lambda.amazonaws.com"
		lambdaFunctionName   = name(e.Project, "lambda-function")
	)

	vpc := i.NewVpc(vpcName, Cidr)
	securityGroup := i.NewSecurityGroup(securityGroupName, vpc)

	i.NewVpcInterfaceEndpoint(resources.NewVpcEndpointProps{
		SecurityGroup: securityGroup,
		ServiceName:   name(e.Project, "ecr"),
		Service:       ec2.InterfaceVpcEndpointAwsService_ECR(),
		Subnets:       vpc.PrivateSubnets(),
		Vpc:           vpc,
	})

	i.NewVpcInterfaceEndpoint(resources.NewVpcEndpointProps{
		SecurityGroup: securityGroup,
		ServiceName:   name(e.Project, "dkr-ecr"),
		Service:       ec2.InterfaceVpcEndpointAwsService_ECR_DOCKER(),
		Subnets:       vpc.PrivateSubnets(),
		Vpc:           vpc,
	})

	i.NewVpcInterfaceEndpoint(resources.NewVpcEndpointProps{
		SecurityGroup: securityGroup,
		ServiceName:   name(e.Project, "logs-ecr"),
		Service:       ec2.InterfaceVpcEndpointAwsService_CLOUDWATCH_LOGS(),
		Subnets:       vpc.PrivateSubnets(),
		Vpc:           vpc,
	})

	logBucket := i.GetBucketFromName(LogBucketName)
	logGroup := i.GetLogGroupFromName(LogGroupName)

	repository := i.NewEcrRepository(repositoryName)
	i.NewCluster(resources.NewClusterProps{
		Name:      clusterName,
		LogBucket: logBucket,
		LogGroup:  logGroup,
		Vpc:       vpc,
	})
	taskDefinition := i.NewTaskDefinition(taskName)
	i.AddContainer(resources.AddContainerProps{
		Name:           containerName,
		Tag:            "latest",
		LogGroup:       logGroup,
		TaskDefinition: taskDefinition,
		Repository:     repository,
	})

	deadLetterQueue := i.NewQueue(deadLetterQueueName, nil)
	queue := i.NewQueue(queueName, &awssqs.DeadLetterQueue{
		MaxReceiveCount: jsii.Number(1),
		Queue:           deadLetterQueue,
	})

	lambdaAssumeRole := i.NewAssumeRole(lambdaAssumeRoleName, lambdaPrincipal)
	i.AddPolicyToRole(resources.AddPolicyToRoleProps{
		Role:      lambdaAssumeRole,
		Actions:   []string{"sqs:*", "logs:*", "ecs:RunTask", "iam:PassRole"},
		Resources: []string{"*"},
	})

	i.NewLambdaFunction(resources.NewLambdaFunctionProps{
		CodePath:       "./bin/handler.zip",
		HandlerPath:    "infra/bin/handler",
		EventSourceArn: queue.QueueArn(),
		Name:           lambdaFunctionName,
		Role:           lambdaAssumeRole,
		Env: map[string]string{
			"Cluster":        clusterName,
			"Task":           *taskDefinition.TaskDefinitionArn(),
			"Container":      containerName,
			"Subnets":        i.RetrieveSubnetIds(vpc),
			"SecurityGroups": *securityGroup.SecurityGroupId(),
		},
	})

	return stack
}

const (
	BootstrapBucketName string = "BBN"
	ConnectionArn       string = "CARN"
	Env                 string = "ENV"
	GithubAccessToken   string = "GHAT"
	GithubOwner         string = "GHO"
	GithubRepository    string = "GHR"
	HostedZoneId        string = "HGI"
	Id                  string = "ID"
	Project             string = "PROJECT"
)

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)
	var (
		bbn     = app.Node().TryGetContext(jsii.String(BootstrapBucketName))
		carn    = app.Node().TryGetContext(jsii.String(ConnectionArn))
		env     = app.Node().TryGetContext(jsii.String(Env))
		ght     = app.Node().TryGetContext(jsii.String(GithubAccessToken))
		gho     = app.Node().TryGetContext(jsii.String(GithubOwner))
		ghr     = app.Node().TryGetContext(jsii.String(GithubRepository))
		hgi     = app.Node().TryGetContext(jsii.String(HostedZoneId))
		id      = app.Node().TryGetContext(jsii.String(Id))
		project = app.Node().TryGetContext(jsii.String(Project))
	)

	if bbn == nil || carn == nil || env == nil || ght == nil || gho == nil || ghr == nil || hgi == nil || project == nil || id == nil {
		panic("please pass context")
	}

	awscdk.Tags_Of(app).Add(jsii.String("Project"), jsii.String(project.(string)), nil)

	NewInfraStack(app, fmt.Sprintf("%sStack", project.(string)),
		&InfraStackProps{
			awscdk.StackProps{
				Env: myenv(),
				Synthesizer: awscdk.NewDefaultStackSynthesizer(
					&awscdk.DefaultStackSynthesizerProps{
						FileAssetsBucketName: jsii.String(bbn.(string)),
						BucketPrefix:         jsii.String(fmt.Sprintf("%sStack/", project.(string))),
					},
				),
			},
		},
		Props{
			Project: project.(string),
		},
	)

	app.Synth(nil)
}

func myenv() *awscdk.Environment { return nil }
