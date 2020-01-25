package parser

import "testing"

func TestParseHostname(t *testing.T) {

	cases := []struct {
		in, out string
	}{
		{"123.123.123.123", ""},
		{"ip-123-123-123-123", "123.123.123.123"},
		{"ip-123-123-123-123.ec2.internal", "123.123.123.123"},
	}

	for _, c := range cases {
		response := parseHostname(c.in)
		if response != c.out {
			t.Fatalf("Response %s did not match expected output %s", response, c.out)
		}
	}
}
