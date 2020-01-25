package command

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/NoUseFreak/sawsh/internal/sawsh"
	"github.com/NoUseFreak/sawsh/internal/sawsh/parser"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(connectCmd)
	addConnectFlags(connectCmd)
}

func addConnectFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().Bool("ssm", false, "Try connecting using AWS Service Manager")
}

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to an instance",
	Args:  cobra.ExactArgs(1),
	Run:   connect,
}

func connect(cmd *cobra.Command, args []string) {
	debug, _ := cmd.Flags().GetBool("debug")

	r := regexp.MustCompile("(([^@]+)@)?(.+)")
	m := r.FindAllStringSubmatch(args[0], 3)
	opts := sawsh.SSHOptions{
		User: m[0][2],
	}
	if b, _ := cmd.Flags().GetBool("ssm"); b {
		opts.TrySSM = b
	}

	sshConnect(opts, *doResolve(cmd, m[0][3]), debug)
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func doResolve(cmd *cobra.Command, hostname string) *sawsh.Instance {
	for _, r := range parser.GetParsers() {
		if i, err := r(cmd, hostname); i != nil && err == nil {
			logrus.Debugf("Found %s using %s", i, getFunctionName(r))
			return i
		}
	}
	logrus.Debugf("Unable to parse input")
	return &sawsh.Instance{
		Ip: hostname,
	}
}

func sshConnect(opts sawsh.SSHOptions, instance sawsh.Instance, debug bool) error {
	if instance.PublicIp != "" && isPortOpen(instance.PublicIp, 22) {
		return executeSSHConnect(opts, instance.PublicIp, nil, debug)
	}
	if isPortOpen(instance.Ip, 22) {
		return executeSSHConnect(opts, instance.Ip, nil, debug)
	}
	logrus.Debug("Remote port is not open, looking for other options")
	if opts.TrySSM {
		ssmOpts := []string{"-o", `ProxyCommand='/bin/sh -c "aws ssm start-session --target %h --document-name AWS-StartSSHSession --parameters 'portNumber=%p'"'`}
		return executeSSHConnect(opts, instance.InstanceId, ssmOpts, debug)
	}

	logrus.Debug("No other options, trying anyway")
	return executeSSHConnect(opts, instance.Ip, nil, debug)
}

func executeSSHConnect(opts sawsh.SSHOptions, hostname string, options []string, debug bool) error {
	addr := hostname
	if opts.User != "" {
		addr = fmt.Sprintf("%s@%s", opts.User, hostname)
	}
	args := append([]string{"ssh", addr}, options...)

	fmt.Printf("Connecting to %s...\n", addr)

	cmd := exec.Command(os.Getenv("SHELL"), []string{"-c", strings.Join(args, " ")}...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if debug {
		fmt.Println(cmd.String())
	} else {
		cmd.Run()
	}

	return nil
}
func isPortOpen(host string, port int) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, time.Duration(1)*time.Second)
	if conn != nil {
		defer conn.Close()
	}

	if err, ok := err.(*net.OpError); ok && err.Timeout() {
		return false
	}

	if err != nil {
		return false
	}

	return true
}
