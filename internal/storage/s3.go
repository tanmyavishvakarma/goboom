package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"goboom/internal/helper"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func CreateS3Session() (*s3.Client, error) {

	cfg, err := helper.GetAWSConfig()
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(cfg), nil
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
