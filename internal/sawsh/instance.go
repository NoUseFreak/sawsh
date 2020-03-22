package sawsh

import "time"

// Instance is a collection of all information about an ec2 instance.
type Instance struct {
	// Name of the instance
	Name string
	// InstanceId is the AWS ec2 identifier
	InstanceID string
	// Ip is the private IP address.
	IP string
	// PublicIp is the public IP address.
	PublicIP string
	// LaunchTime is the time the ec2 instance was launched
	LaunchTime time.Time
}
