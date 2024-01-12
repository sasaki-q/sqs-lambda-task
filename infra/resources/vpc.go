package resources

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/jsii-runtime-go"
)

type VpcResource struct {
	Vpc              ec2.Vpc
	SubnetIds        string
	SecurityGroupIds string
}

func NewVpc(s awscdk.Stack) VpcResource {
	v := ec2.NewVpc(s, jsii.String("vpc"), &ec2.VpcProps{
		VpcName: jsii.String("vpc"),
		Cidr:    jsii.String("192.168.0.0/16"),
		SubnetConfiguration: &[]*ec2.SubnetConfiguration{
			{Name: jsii.String("public"), CidrMask: jsii.Number(24), SubnetType: ec2.SubnetType_PUBLIC},
			{Name: jsii.String("private"), CidrMask: jsii.Number(24), SubnetType: ec2.SubnetType_PRIVATE_ISOLATED},
		},
	})

	v.IsolatedSubnets()

	v.AddGatewayEndpoint(jsii.String("s3-endpoint"),
		&ec2.GatewayVpcEndpointOptions{
			Service: ec2.GatewayVpcEndpointAwsService_S3(),
			Subnets: &[]*ec2.SubnetSelection{{SubnetType: ec2.SubnetType_PRIVATE_ISOLATED}},
		},
	)

	ecrep := v.AddInterfaceEndpoint(jsii.String("ecr-endpoint"),
		&ec2.InterfaceVpcEndpointOptions{
			Service:           ec2.InterfaceVpcEndpointAwsService_ECR(),
			PrivateDnsEnabled: jsii.Bool(true),
			Subnets:           &ec2.SubnetSelection{SubnetType: ec2.SubnetType_PRIVATE_ISOLATED},
		},
	)

	ecrdkrep := v.AddInterfaceEndpoint(jsii.String("ecr-dkr-endpoint"),
		&ec2.InterfaceVpcEndpointOptions{
			Service:           ec2.InterfaceVpcEndpointAwsService_ECR_DOCKER(),
			PrivateDnsEnabled: jsii.Bool(true),
			Subnets:           &ec2.SubnetSelection{SubnetType: ec2.SubnetType_PRIVATE_ISOLATED},
		},
	)

	cwlep := v.AddInterfaceEndpoint(jsii.String("cloudwatch-logs-endpoint"),
		&ec2.InterfaceVpcEndpointOptions{
			Service:           ec2.InterfaceVpcEndpointAwsService_CLOUDWATCH_LOGS(),
			PrivateDnsEnabled: jsii.Bool(true),
			Subnets:           &ec2.SubnetSelection{SubnetType: ec2.SubnetType_PRIVATE_ISOLATED},
		},
	)

	return VpcResource{
		Vpc:              v,
		SubnetIds:        retrieveSubnetIds(v),
		SecurityGroupIds: fmt.Sprintf("%s,%s,%s", retrieveSgIds(ecrep), retrieveSgIds(ecrdkrep), retrieveSgIds(cwlep)),
	}
}

func retrieveSubnetIds(v ec2.Vpc) string {
	subnetIds := ""
	for i, v := range *v.IsolatedSubnets() {
		if i == 0 {
			subnetIds += *v.SubnetId()
		} else {
			subnetIds += fmt.Sprintf(",%s", *v.SubnetId())
		}
	}

	return subnetIds
}

func retrieveSgIds(ep ec2.InterfaceVpcEndpoint) string {
	sg := ""
	for i, v := range *ep.Connections().SecurityGroups() {
		if i == 0 {
			sg += *v.SecurityGroupId()
		} else {
			sg += *v.SecurityGroupId()
		}
	}
	return sg
}
