package aws

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func RegionCode(region string) string {
	return strings.ReplaceAll(strings.ToLower(region), "_", "-")
}

func newSession(region string) (*session.Session, error) {
	awsSession, err := session.NewSession(aws.NewConfig().WithRegion(region))
	if err != nil {
		return nil, err
	}
	return awsSession, nil
}
