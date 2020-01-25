package command

import (
	"github.com/NoUseFreak/sawsh/internal/sawsh/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.PersistentFlags().Bool("plain", false, "Print only the resulting ip's")
	listCmd.PersistentFlags().Bool("csv", false, "Print ip's in csv format")
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List instances",
	Args:  cobra.ExactArgs(1),
	Run:   list,
}

func list(cmd *cobra.Command, args []string) {
	instances := utils.ListInstances(args[0])

	if b, _ := cmd.Flags().GetBool("plain"); b {
		utils.PrintPlain(instances)
	} else if b, _ := cmd.Flags().GetBool("csv"); b {
		utils.PrintCSV(instances)
	} else {
		utils.PrintTable(instances)
	}
}
