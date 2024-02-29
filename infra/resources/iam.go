package resources

import (
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/jsii-runtime-go"
)

func (r *ResourceService) NewAssumeRole(name string, principal string) iam.Role {
	return iam.NewRole(r.S, jsii.String(name), &iam.RoleProps{
		AssumedBy: iam.NewServicePrincipal(jsii.String(principal), nil),
		RoleName:  jsii.String(name),
	})
}

func (r *ResourceService) AddPolicyToRole(e AddPolicyToRoleProps) {
	e.Role.AddToPolicy(iam.NewPolicyStatement(&iam.PolicyStatementProps{
		Actions:   sVtoP(e.Actions),
		Effect:    iam.Effect_ALLOW,
		Resources: sVtoP(e.Resources),
	}))
}
