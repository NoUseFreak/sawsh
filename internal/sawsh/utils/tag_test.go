package utils

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestFindTag(t *testing.T) {
	cases := []struct {
		key, value string
	}{
		{"Name", "Name Value"},
		{"Environment", "Production"},
		{"Fake", ""},
	}
	tags := []*ec2.Tag{
		{
			Key:   aws.String("Name"),
			Value: aws.String("Name Value"),
		},
		{
			Key:   aws.String("Environment"),
			Value: aws.String("Production"),
		},
	}

	for _, c := range cases {
		value := findTag(c.key, tags)

		if value != c.value {
			t.Fatalf("Could not resolve Name")
		}
	}
}
