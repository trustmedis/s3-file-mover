package lib

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func CreateSession(config *Config) (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region:   aws.String(config.REGION),
		Endpoint: aws.String(config.ENDPOINT),
		Credentials: credentials.NewStaticCredentials(
			config.ACCESS_KEY,
			config.ACCESS_SECRET,
			"",
		),
	})
}

func UploadFile(config *Config, originFilePath, targetFilePath string) error {
	session, err := CreateSession(config)
	if err != nil {
		return nil
	}
	svc := s3.New(session)

	file, err := os.Open(originFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	if config.APPEND_TIMESTAMP {
		targetFilePath = fileInfo.ModTime().Format("20060102150405") + "_" + targetFilePath
	}
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(config.BUCKET),
		Key:           aws.String(targetFilePath),
		Body:          file,
		ContentLength: aws.Int64(fileInfo.Size()),
	})
	if err != nil {
		return err
	}

	log.Printf("Found file %s, uploaded to s3://%s/%s", originFilePath, config.BUCKET, targetFilePath)

	// Autocleanup based on AUTO_CLEANUP value
	if config.AUTO_CLEANUP {
		os.Remove(originFilePath)
	}

	return nil
}
