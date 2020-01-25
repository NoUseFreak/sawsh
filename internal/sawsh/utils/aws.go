package utils

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/NoUseFreak/sawsh/internal/sawsh"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func ListInstances(hostname string) []sawsh.Instance {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	svc := ec2.New(sess)

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(strings.Join([]string{"*", hostname, "*"}, "")),
				},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running"), aws.String("pending")},
			},
		},
	}
	resp, err := svc.DescribeInstances(params)
	if err != nil {
		fmt.Println("there was an error listing instances", err.Error())
		log.Fatal(err.Error())
	}

	var instances []sawsh.Instance
	for _, reservation := range resp.Reservations {
		for _, i := range reservation.Instances {
			inst := sawsh.Instance{
				Name:       findTag("Name", i.Tags),
				Ip:         *i.PrivateIpAddress,
				InstanceId: *i.InstanceId,
				LaunchTime: *i.LaunchTime,
			}
			if i.PublicIpAddress != nil {
				inst.PublicIp = *i.PublicIpAddress
			}
			instances = append(instances, inst)

		}
	}

	// Sort by name, ip
	sort.Slice(instances, func(i, j int) bool {
		if instances[i].Name == instances[j].Name {
			return instances[i].Ip < instances[j].Ip
		}
		return instances[i].Name < instances[j].Name
	})

	return instances
}
