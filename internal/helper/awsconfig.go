package helper

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func GetAWSConfig() (aws.Config, error) {

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),

		config.WithRegion(
			os.Getenv("AWS_REGION"),
		),
	)
	if err != nil {
		return aws.Config{}, err
	}

	return cfg, nil
}


