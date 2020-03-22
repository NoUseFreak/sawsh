package parser

import (
	"github.com/NoUseFreak/sawsh/internal/sawsh"
	"github.com/spf13/cobra"
)

// Parser is an interface to loopup instance information.
type Parser func(cmd *cobra.Command, hostname string) (*sawsh.Instance, error)

// GetParsers returns a list of Parsers.
func GetParsers() []Parser {
	return []Parser{
		ipParser,
		hostnameParser,
		lookupParser,
	}
}
