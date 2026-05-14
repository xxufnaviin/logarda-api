package utils

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func VerifyAWSCredentials(accessKey, secretKey, region string) error {
	cfg := aws.Config{
		Credentials: credentials.NewStaticCredentialsProvider(accessKey,secretKey,"",),
		Region: region,
	}
	client := sts.NewFromConfig(cfg)
	_, err := client.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})

	return err
}
