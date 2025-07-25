package s3

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	aws_s3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

type Config struct {
	Region string
	EndpointURL  string
	AccessKeyID  string
	SecretAccessKey string
}

func Connect(opts Config) (*aws_s3.Client) {
	return aws_s3.New(aws_s3.Options{
		 AppID: "my-application/0.0.1",
		 
		 Region: opts.Region,
		 BaseEndpoint: aws.String(opts.EndpointURL),

		 Credentials: credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     opts.AccessKeyID,
				SecretAccessKey: opts.SecretAccessKey,
			}},
	})
}
