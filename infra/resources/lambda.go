package resources

import (
	lambda "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/jsii-runtime-go"
)

func (r *ResourceService) NewLambdaFunction(e NewLambdaFunctionProps) lambda.Function {
	function := lambda.NewFunction(r.S, jsii.String(e.Name), &lambda.FunctionProps{
		FunctionName: jsii.String(e.Name),
		Handler:      jsii.String(e.HandlerPath),
		Runtime:      lambda.Runtime_GO_1_X(),
		Code:         lambda.AssetCode_FromAsset(jsii.String(e.CodePath), nil),
		Role:         e.Role,
		Environment:  mVtoP(e.Env),
	})

	if e.EventSourceArn != nil {
		function.AddEventSourceMapping(jsii.String("es-mapping"), &lambda.EventSourceMappingOptions{
			EventSourceArn: e.EventSourceArn,
		})
	}

	return function
}
