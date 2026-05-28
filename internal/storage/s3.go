package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func CreateS3Session() (*s3.Client, error) {

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),

		config.WithRegion(
			os.Getenv("AWS_REGION"),
		),
	)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg)

	return s3Client, nil
}

func UploadDirectoryToS3(bucketName string, deploymentID string, directory string) error {
	client, err := CreateS3Session()
	if err != nil {
		return err
	}

	return filepath.WalkDir(directory, func(filePath string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		relativePath, err := filepath.Rel(directory, filePath)
		if err != nil {
			return err
		}

		objectKey := fmt.Sprintf("%s/%s", deploymentID, relativePath)

		return uploadFileToS3(client, bucketName, objectKey, filePath)
	})
}

func uploadFileToS3(client *s3.Client, bucketName string, objectKey string, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
		Body:   file,
	})
	defer file.Close()
	return err
}
