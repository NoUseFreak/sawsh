package parser

import (
	"fmt"
	"regexp"

	"github.com/NoUseFreak/sawsh/internal/sawsh"
	"github.com/spf13/cobra"
)

func hostnameParser(cmd *cobra.Command, hostname string) (*sawsh.Instance, error) {
	h := parseHostname(hostname)
	if h == "" {
		return nil, nil
	}

	return &sawsh.Instance{
		IP: h,
	}, nil
}

func parseHostname(hostname string) string {
	r := regexp.MustCompile("ip-([0-9]{1,3})-([0-9]{1,3})-([0-9]{1,3})-([0-9]{1,3})")
	m := r.FindAllStringSubmatch(hostname, 4)

	if len(m) != 1 {
		return ""
	}
	return fmt.Sprintf("%s.%s.%s.%s", m[0][1], m[0][2], m[0][3], m[0][4])
}
