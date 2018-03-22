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
	"time"

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

	regionFlag := cli.StringFlag{
		Name:  "aws-region",
		Value: "us-east-1",
	}

	app.Flags = []cli.Flag{
		regionFlag,
	}

	app.Commands = []cli.Command{
		{
			Name:      "connect",
			Usage:     "Search and connect to an instance",
			ArgsUsage: "[filter]",
			Flags: []cli.Flag{
				regionFlag,
				cli.BoolFlag{
					Name:  "transparant",
					Usage: "Connect if the ec2 query has no results",
				},
				cli.BoolFlag{
					Name:  "print",
					Usage: "Print instead of connecting",
				},
			},
			Action: connectAction,
		},
		{
			Name:      "exec",
			Usage:     "Search and execute a command on multiple servers",
			ArgsUsage: "[filter] [command]",
			Flags: []cli.Flag{
				regionFlag,
			},
			Action: execAction,
		},
		{
			Name:   "list",
			Usage:  "Render a list of instances",
			Action: listAction,
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func connectAction(c *cli.Context) error {
	hostname := c.Args().First()
	instances, err := findInstances(c, hostname)
	if err != nil && c.Bool("transparant") {
		return sshConnect(c, hostname)
	} else if len(instances) == 0 {
		return nil
	}
	var instance Instance
	if len(instances) > 1 {
		printTable(instances)
		instance = getTableChoice(instances)
	} else {
		instance = instances[0]
	}
	return sshConnect(c, instance.ip)
}

func execAction(c *cli.Context) error {
	hostname := c.Args().First()
	command := c.Args().Get(1)
	instances, _ := findInstances(c, hostname)
	ips := getInstanceIps(instances)

	parallelCmd := []string{
		"parallel",
		"--no-notice",
		"--tagstring",
		"{}",
		fmt.Sprintf("ssh {} %s", command),
		":::",
	}
	parallelCmd = append(parallelCmd, ips...)
	//fmt.Println(strings.Join(parallelCmd, " "))

	cmd := exec.Command(parallelCmd[0], parallelCmd[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()

	return nil
}

func listAction(c *cli.Context) error {
	hostname := c.Args().First()
	instances := queryAws(hostname, "us-east-1")
	printTable(instances)
	return nil
}

// Find the hostname's ip
func findInstances(c *cli.Context, hostname string) ([]Instance, error) {
	var err error
	ip, err := parseHostname(hostname)
	if err == nil {
		return []Instance{Instance{ip: ip}}, err
	}
	instances := queryAws(hostname, c.String("aws-region"))
	if len(instances) == 0 {
		fmt.Println("No instances found")
		return instances, errors.New("No instance found")
	}
	return instances, err
}

// Check if the hostname is a known aws hostname format `ip-123-123-123-123` or just a plain ip.
func parseHostname(hostname string) (string, error) {
	var err error
	if net.ParseIP(hostname) != nil {
		return hostname, err
	}
	r, _ := regexp.Compile("ip-([0-9]{1,3})-([0-9]{1,3})-([0-9]{1,3})-([0-9]{1,3})")
	m := r.FindAllStringSubmatch(hostname, 4)

	if len(m) != 1 {
		return "", errors.New("Unknown hostname")
	}

	return fmt.Sprintf("%s.%s.%s.%s", m[0][1], m[0][2], m[0][3], m[0][4]), err
}

// Open an interactive ssh connection
func sshConnect(c *cli.Context, hostname string) error {
	if c.Bool("print") {
		fmt.Printf("ssh %s", hostname)
		return nil
	}
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
			instances = append(instances, Instance{
				name:       findTag("Name", instance.Tags),
				ip:         *instance.PrivateIpAddress,
				launchTime: *instance.LaunchTime,
			})
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
	table.SetHeader([]string{"", "Name", "Ip", "LaunchTime"})

	for index, instance := range instances {
		table.Append([]string{strconv.Itoa(index), instance.name, instance.ip, instance.launchTime.String()})
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
	name       string
	ip         string
	launchTime time.Time
}

func getInstanceIps(instances []Instance) []string {
	vsm := make([]string, len(instances))
	for i, v := range instances {
		vsm[i] = v.ip
	}
	return vsm
}
