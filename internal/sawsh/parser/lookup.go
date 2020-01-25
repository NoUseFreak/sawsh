package parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/NoUseFreak/sawsh/internal/sawsh"
	"github.com/NoUseFreak/sawsh/internal/sawsh/utils"
	"github.com/spf13/cobra"
)

func lookupParser(cmd *cobra.Command, hostname string) (*sawsh.Instance, error) {
	instances := utils.ListInstances(hostname)

	switch len(instances) {
	case 0:
		return nil, nil
	case 1:
		return &instances[0], nil
	}

	utils.PrintTable(instances)
	i := getTableChoice(instances)
	return &i, nil
}

func getTableChoice(instances []sawsh.Instance) sawsh.Instance {
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

type instance struct {
	name       string
	ip         string
	launchTime time.Time
}
