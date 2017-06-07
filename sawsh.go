package main

import (
	"bufio"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olekukonko/tablewriter"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	arg := os.Args[1]

	var ip string
	if net.ParseIP(arg) != nil {
		ip = arg
	} else {
		ip = queryEC2(arg)
	}

	fmt.Println("Connecting to", ip, "...\n")

	cmd := exec.Command("ssh", ip)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func queryEC2(filter string) string {
	awsRegion := "us-east-1"

	instances := findInstances(filter, awsRegion)
	if len(instances) == 1 {
		return instances[0].ip
	} else if len(instances) == 0 {
		fmt.Println("Could not match any instances")
		os.Exit(1)
	}
	fmt.Printf("listing instances with tag %v in: %v\n", filter, awsRegion)
	printTable(instances)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Pick a number: ")
	choice, _ := reader.ReadString('\n')
	choiceInt, err := strconv.Atoi(choice[:len(choice)-1])
	if err != nil {
		choiceInt = 0
	}

	return instances[choiceInt].ip
}

func printTable(instances []Instance) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"", "Name", "Ip"})

	for index, instance := range instances {
		table.Append([]string{strconv.Itoa(index), instance.name, instance.ip})
	}

	table.Render()
}

func findInstances(filter string, awsRegion string) []Instance {
	sess := session.Must(session.NewSession())

	svc := ec2.New(sess, &aws.Config{Region: aws.String(awsRegion)})
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(strings.Join([]string{"*", filter, "*"}, "")),
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
		fmt.Println("there was an error listing instances in", awsRegion, err.Error())
		log.Fatal(err.Error())
	}

	var instances []Instance
	for _, instance := range resp.Reservations {
		i := instance.Instances[0]
		instances = append(instances, Instance{name: findTag("Name", i.Tags), ip: *i.PrivateIpAddress})
	}
	return instances
}

func findTag(key string, tags []*ec2.Tag) string {

	for _, item := range tags {
		if *item.Key == key {
			return *item.Value
		}
	}

	return ""
}

type Instance struct {
	name string
	ip   string
}
