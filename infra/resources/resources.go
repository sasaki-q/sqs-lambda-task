package resources

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	ecr "github.com/aws/aws-cdk-go/awscdk/v2/awsecr"
	ecs "github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	lambda "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	logs "github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	s3 "github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	sqs "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
)

type ResourceService struct {
	S cdk.Stack
}

type IResourceService interface {
	// ecr.go
	NewEcrRepository(name string) ecr.Repository
	GetEcrRepositoryFromName(name string) ecr.IRepository

	// ecs.go
	NewCluster(e NewClusterProps) ecs.Cluster
	NewTaskDefinition(name string) ecs.TaskDefinition
	AddContainer(e AddContainerProps) ecs.ContainerDefinition

	// iam.go
	NewAssumeRole(name string, principal string) iam.Role
	AddPolicyToRole(e AddPolicyToRoleProps)

	// lambda.go
	NewLambdaFunction(e NewLambdaFunctionProps) lambda.Function

	// logs.go
	GetLogGroupFromName(name string) logs.ILogGroup

	// s3.go
	NewBucket(name string) s3.Bucket
	GetBucketFromName(name string) s3.IBucket

	// sqs.go
	NewQueue(name string, dq *sqs.DeadLetterQueue) sqs.Queue

	// vpc.go
	NewVpc(vpcName string, cidr string) ec2.Vpc
	NewSecurityGroup(name string, vpc ec2.IVpc) ec2.SecurityGroup
	NewVpcInterfaceEndpoint(e NewVpcEndpointProps) ec2.InterfaceVpcEndpoint
	RetrieveSubnetIds(vpc ec2.Vpc) string
}

type NewClusterProps struct {
	Name      string
	LogBucket s3.IBucket
	LogGroup  logs.ILogGroup
	Vpc       ec2.Vpc
}

type AddContainerProps struct {
	Name string
	Tag  string

	LogGroup       logs.ILogGroup
	TaskDefinition ecs.TaskDefinition
	Repository     ecr.IRepository
}

type AddPolicyToRoleProps struct {
	Actions   []string
	Resources []string
	Role      iam.Role
}

type NewLambdaFunctionProps struct {
	Name           string
	CodePath       string
	HandlerPath    string
	Env            map[string]string
	EventSourceArn *string
	Role           iam.Role
}

type NewVpcEndpointProps struct {
	SecurityGroup ec2.SecurityGroup
	ServiceName   string
	Service       ec2.IInterfaceVpcEndpointService
	Subnets       *[]ec2.ISubnet
	Vpc           ec2.Vpc
}
