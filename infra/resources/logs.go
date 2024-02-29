package resources

import (
	logs "github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/jsii-runtime-go"
)

func (r *ResourceService) GetLogGroupFromName(name string) logs.ILogGroup {
	return logs.LogGroup_FromLogGroupName(r.S, jsii.String(name), jsii.String(name))
}
