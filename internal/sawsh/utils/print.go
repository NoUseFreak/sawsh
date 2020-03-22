package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/NoUseFreak/sawsh/internal/sawsh"
	"github.com/olekukonko/tablewriter"
)

// PrintTable prints a list of instances in table format.
func PrintTable(instances []sawsh.Instance) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"", "Name", "Ip", "LaunchTime"})

	for index, i := range instances {
		table.Append([]string{strconv.Itoa(index), i.Name, i.IP, i.LaunchTime.String()})
	}
	table.Render()
}

// PrintPlain prints a list of instances al a list of ip's.
func PrintPlain(instances []sawsh.Instance) {
	for _, instance := range instances {
		fmt.Println(instance.IP)
	}
}

// PrintCSV prints a list of instances in csv format.
func PrintCSV(instances []sawsh.Instance) {
	ips := []string{}
	for _, instance := range instances {
		ips = append(ips, instance.IP)
	}
	fmt.Println(strings.Join(ips, ","))
}
