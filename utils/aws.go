package utils

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func VerifyAWSCredentials(accessKey, secretKey, region string) error {
	cfg := aws.Config{
		Credentials: credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		Region:      region,
	}
	client := sts.NewFromConfig(cfg)
	_, err := client.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})

	return err
}

func HasAllPermissions(accessKey, secretKey, region string) bool {
	ctx := context.Background()
	cfg := aws.Config{
		Credentials: credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		Region:      region,
	}

	ec2Client := ec2.NewFromConfig(cfg)
	if _, err := ec2Client.DescribeRegions(ctx, &ec2.DescribeRegionsInput{}); err != nil {
		return false
	}

	cwClient := cloudwatch.NewFromConfig(cfg)
	if _, err := cwClient.ListMetrics(ctx, &cloudwatch.ListMetricsInput{}); err != nil {
		return false
	}

	ctClient := cloudtrail.NewFromConfig(cfg)
	if _, err := ctClient.LookupEvents(ctx, &cloudtrail.LookupEventsInput{}); err != nil {
		return false
	}

	return true
}
