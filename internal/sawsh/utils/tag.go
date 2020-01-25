package utils

import "github.com/aws/aws-sdk-go/service/ec2"

func findTag(key string, tags []*ec2.Tag) string {

	for _, item := range tags {
		if *item.Key == key {
			return *item.Value
		}
	}

	return ""
}
