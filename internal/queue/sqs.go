package queue

import (
	"goboom/internal/helper"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func CreateSQSSession() (*sqs.Client, error) {
	cfg, err := helper.GetAWSConfig()
	if err != nil {
		return nil, err
	}

	return sqs.NewFromConfig(cfg), nil
}
