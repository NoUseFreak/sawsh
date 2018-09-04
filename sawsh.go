package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"

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
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "plain",
					Usage: "Print only the resulting ip's",
				},
			},
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
	if c.Bool("plain") {
		printPlain(instances)
	} else {
		printTable(instances)
	}

	return nil
}
