package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"testing"
)

func TestParseInput(t *testing.T) {

	cases := []struct {
		in, out string
	}{
		{"123.123.123.123", "123.123.123.123"},
		{"ip-123-123-123-123", "123.123.123.123"},
	}

	for _, c := range cases {
		response, err := parseInput(c.in)
		if err != nil {
			t.Fatal(err)
			t.Fatalf("Failed to parse input %v", err)
		}
		if response != c.out {
			t.Fatalf("Response did not match expeceted output")
		}
	}
}

func TestFindTag(t *testing.T) {
	cases := []struct {
		key, value string
	}{
		{"Name", "Name Value"},
		{"Environment", "Production"},
	}
	tags := []*ec2.Tag{
		&ec2.Tag{
			Key:   aws.String("Name"),
			Value: aws.String("Name Value"),
		},
		&ec2.Tag{
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
