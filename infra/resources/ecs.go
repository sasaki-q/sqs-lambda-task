package resources

import (
	ecs "github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/jsii-runtime-go"
)

func (r *ResourceService) NewCluster(e NewClusterProps) ecs.Cluster {
	return ecs.NewCluster(r.S, jsii.String(e.Name), &ecs.ClusterProps{
		ClusterName: jsii.String(e.Name),
		ExecuteCommandConfiguration: &ecs.ExecuteCommandConfiguration{
			LogConfiguration: &ecs.ExecuteCommandLogConfiguration{
				CloudWatchLogGroup:          e.LogGroup,
				CloudWatchEncryptionEnabled: jsii.Bool(true),
				S3Bucket:                    e.LogBucket,
				S3EncryptionEnabled:         jsii.Bool(true),
				S3KeyPrefix:                 jsii.String(e.Name),
			},
			Logging: ecs.ExecuteCommandLogging_OVERRIDE,
		},
		Vpc: e.Vpc,
	})
}

func (r *ResourceService) NewTaskDefinition(name string) ecs.TaskDefinition {
	return ecs.NewFargateTaskDefinition(r.S, jsii.String(name), &ecs.FargateTaskDefinitionProps{
		Cpu:             jsii.Number(256),
		MemoryLimitMiB:  jsii.Number(512),
		Family:          jsii.String(name),
		RuntimePlatform: &ecs.RuntimePlatform{CpuArchitecture: ecs.CpuArchitecture_X86_64()},
	})
}

func (r *ResourceService) AddContainer(e AddContainerProps) ecs.ContainerDefinition {
	return e.TaskDefinition.AddContainer(jsii.String(e.Name), &ecs.ContainerDefinitionOptions{
		ContainerName: jsii.String(e.Name),
		Image:         ecs.ContainerImage_FromEcrRepository(e.Repository, jsii.String(e.Tag)),
		Logging: ecs.AwsLogDriver_AwsLogs(&ecs.AwsLogDriverProps{
			StreamPrefix: jsii.String(e.Name),
			LogGroup:     e.LogGroup,
		}),
	})
}
