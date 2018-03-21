package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "sawsh"
	app.Usage = "Query and connect to ec2 instances"
	app.HideVersion = true
	app.Copyright = "(c) Dries De Peuter <dries@depeuter.io>"
	app.ArgsUsage = "[filter]"
	app.Action = connectAction
	app.EnableBashCompletion = true

	app.Commands = []cli.Command{
		{
			Name:   "list",
			Usage:  "Render a list of instances",
			Action: listAction,
		},
		{
			Name:      "connect",
			Usage:     "Search and connect to an instance",
			ArgsUsage: "[filter]",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "transparant",
					Usage: "Connect if the ec2 query has no results",
				},
			},
			Action: connectAction,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func connectAction(c *cli.Context) error {
	hostname := c.Args().First()
	ip, err := parseHostname(hostname)
	if err == nil {
		return sshConnect(ip)
	}
	instances := queryAws(hostname, "us-east-1")
	if len(instances) == 0 {
		fmt.Println("No instances found")

		if c.Bool("transparant") {
			sshConnect(hostname)
		}
		return nil
	} else if len(instances) == 1 {
		sshConnect(instances[0].ip)
		return nil
	}
	printTable(instances)
	instance := getTableChoice(instances)
	sshConnect(instance.ip)
	return nil
}

func listAction(c *cli.Context) error {
	hostname := c.Args().First()
	instances := queryAws(hostname, "us-east-1")
	printTable(instances)
	return nil
}

// Check if the hostname is a known aws hostname format `ip-123-123-123-123` or just a plain ip.
func parseHostname(hostname string) (string, error) {
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

// Open an interactive ssh connection
func sshConnect(hostname string) error {
	fmt.Printf("Connecting to %s...\n", hostname)
	cmd := exec.Command("ssh", hostname)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()

	return nil
}

func queryAws(filter string, awsRegion string) []Instance {
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
	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			instances = append(instances, Instance{name: findTag("Name", instance.Tags), ip: *instance.PrivateIpAddress})
		}
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

func printTable(instances []Instance) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"", "Name", "Ip"})

	for index, instance := range instances {
		table.Append([]string{strconv.Itoa(index), instance.name, instance.ip})
	}

	table.Render()
}

func getTableChoice(instances []Instance) Instance {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Pick a number: ")
	choice, _ := reader.ReadString('\n')
	choiceInt, err := strconv.Atoi(choice[:len(choice)-1])
	if err != nil {
		choiceInt = 0
	}

	if choiceInt < 0 || choiceInt > len(instances) {
		choiceInt = 0
	}

	return instances[choiceInt]
}

type Instance struct {
	name string
	ip   string
}
