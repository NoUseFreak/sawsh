package sawsh

// SSHOptions contains ssh connection options.
type SSHOptions struct {
	User     string
	Instance Instance
	TrySSM   bool
}
