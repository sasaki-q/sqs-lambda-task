package resources

import (
	sqs "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/jsii-runtime-go"
)

func (r *ResourceService) NewQueue(name string, dq *sqs.DeadLetterQueue) sqs.Queue {
	return sqs.NewQueue(r.S, jsii.String(name), &sqs.QueueProps{
		QueueName:       jsii.String(name),
		DeadLetterQueue: dq,
	})
}
