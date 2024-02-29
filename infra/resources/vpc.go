package resources

import (
	"fmt"
	"strings"

	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/jsii-runtime-go"
)

func (r *ResourceService) NewVpc(vpcName string, cidr string) ec2.Vpc {
	vpc := ec2.NewVpc(r.S, jsii.String(vpcName), &ec2.VpcProps{
		IpAddresses: ec2.IpAddresses_Cidr(jsii.String(cidr)),
		VpcName:     jsii.String(vpcName),
		SubnetConfiguration: &[]*ec2.SubnetConfiguration{
			{Name: jsii.String(fmt.Sprintf("%s-public-", vpcName)), CidrMask: jsii.Number(24), SubnetType: ec2.SubnetType_PUBLIC},
			{Name: jsii.String(fmt.Sprintf("%s-private-", vpcName)), CidrMask: jsii.Number(24), SubnetType: ec2.SubnetType_PRIVATE_WITH_EGRESS},
		},
	})

	vpc.AddGatewayEndpoint(jsii.String("s3-gateway-endpoint"),
		&ec2.GatewayVpcEndpointOptions{
			Service: ec2.GatewayVpcEndpointAwsService_S3(),
			Subnets: &[]*ec2.SubnetSelection{{SubnetType: ec2.SubnetType_PRIVATE_WITH_EGRESS}},
		},
	)

	return vpc
}

func (r *ResourceService) NewSecurityGroup(name string, vpc ec2.IVpc) ec2.SecurityGroup {
	return ec2.NewSecurityGroup(r.S, jsii.String(name), &ec2.SecurityGroupProps{
		SecurityGroupName: jsii.String(name),
		Vpc:               vpc,
	})
}

func (r *ResourceService) NewVpcInterfaceEndpoint(e NewVpcEndpointProps) ec2.InterfaceVpcEndpoint {
	return e.Vpc.AddInterfaceEndpoint(jsii.String(fmt.Sprintf("%s-endpoint", e.ServiceName)),
		&ec2.InterfaceVpcEndpointOptions{
			Service:        e.Service,
			SecurityGroups: &[]ec2.ISecurityGroup{e.SecurityGroup},
			Subnets:        &ec2.SubnetSelection{Subnets: e.Subnets},
		},
	)
}

func (r *ResourceService) RetrieveSubnetIds(vpc ec2.Vpc) string {
	tmp := []string{}
	for _, v := range *vpc.PrivateSubnets() {
		tmp = append(tmp, *v.SubnetId())
	}

	return strings.Join(tmp, ",")
}
