package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olekukonko/tablewriter"
	"log"
	"net"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var config Config

func main() {

	config.transparant = *flag.Bool("transparant", true, "Forward lookup failures")
	flag.Parse()
	if flag.NArg() < 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	config.hostname = flag.Args()[0]

	var ip string
	if hostIp, err := parseInput(config.hostname); err == nil {
		ip = hostIp
	} else {
		ip = queryEC2(config.hostname)
	}

	fmt.Println("Connecting to", ip, "...\n")

	cmd := exec.Command("ssh", ip)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func parseInput(hostname string) (string, error) {
	if net.ParseIP(hostname) != nil {
		return hostname, nil
	}
	var err error
	r, _ := regexp.Compile("ip-([0-9]{1,3})-([0-9]{1,3})-([0-9]{1,3})-([0-9]{1,3})")
	m := r.FindAllStringSubmatch(hostname, 4)

	if len(m) != 1 {
		return "", errors.New("")
	}

	return fmt.Sprintf("%s.%s.%s.%s", m[0][1], m[0][2], m[0][3], m[0][4]), err
}

func queryEC2(filter string) string {
	awsRegion := "us-east-1"

	instances := findInstances(filter, awsRegion)
	if len(instances) == 1 {
		return instances[0].ip
	} else if len(instances) == 0 {
		if config.transparant {
			return filter
		}
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

	// Sort by name, ip
	sort.Slice(instances, func(i, j int) bool {
		if instances[i].name == instances[j].name {
			return instances[i].ip < instances[j].ip
		}
		return instances[i].name < instances[j].name
	})

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

type Target struct {
	ip   string
	user string
}

type Config struct {
	transparant bool
	hostname    string
}
