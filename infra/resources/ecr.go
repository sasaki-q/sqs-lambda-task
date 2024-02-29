package resources

import (
	ecr "github.com/aws/aws-cdk-go/awscdk/v2/awsecr"
	"github.com/aws/jsii-runtime-go"
)

func (r *ResourceService) NewEcrRepository(name string) ecr.Repository {
	return ecr.NewRepository(r.S, jsii.String(name), &ecr.RepositoryProps{
		RepositoryName:     jsii.String(name),
		ImageTagMutability: ecr.TagMutability_IMMUTABLE,
	})
}

func (r *ResourceService) GetEcrRepositoryFromName(name string) ecr.IRepository {
	return ecr.Repository_FromRepositoryName(r.S, jsii.String(name), jsii.String(name))
}
