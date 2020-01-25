package parser

import (
	"github.com/NoUseFreak/sawsh/internal/sawsh"
	"github.com/spf13/cobra"
)

type Parser func(cmd *cobra.Command, hostname string) (*sawsh.Instance, error)

func GetParsers() []Parser {
	return []Parser{
		ipParser,
		hostnameParser,
		lookupParser,
	}
}
