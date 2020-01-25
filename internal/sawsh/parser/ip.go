package parser

import (
	"net"

	"github.com/NoUseFreak/sawsh/internal/sawsh"
	"github.com/spf13/cobra"
)

func ipParser(cmd *cobra.Command, hostname string) (*sawsh.Instance, error) {
	if net.ParseIP(hostname) != nil {
		return &sawsh.Instance{
			Ip: hostname,
		}, nil
	}

	return nil, nil
}
