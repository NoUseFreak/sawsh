package sawsh

import "time"

type Instance struct {
	Name       string
	InstanceId string
	Ip         string
	PublicIp   string
	LaunchTime time.Time
}
