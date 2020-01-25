package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/NoUseFreak/sawsh/internal/sawsh"
	"github.com/olekukonko/tablewriter"
)

func PrintTable(instances []sawsh.Instance) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"", "Name", "Ip", "LaunchTime"})

	for index, i := range instances {
		table.Append([]string{strconv.Itoa(index), i.Name, i.Ip, i.LaunchTime.String()})
	}
	table.Render()
}

func PrintPlain(instances []sawsh.Instance) {
	for _, instance := range instances {
		fmt.Println(instance.Ip)
	}
}

func PrintCSV(instances []sawsh.Instance) {
	ips := []string{}
	for _, instance := range instances {
		ips = append(ips, instance.Ip)
	}
	fmt.Println(strings.Join(ips, ","))
}
