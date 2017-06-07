# sawsh

SSH wrapper for aws to make your life easier.

## Usage

```
sawsh webserver
```

## Install

```
git clone https://github.com/NoUseFreak/sawsh.git
cd sawsh
make install
```
## Configure

Setup sure your [aws-cli](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html) is setup. That user needs `ec2:Describe*` permissions.

### Suggestion

It may be useful to setup some aliases if you use multiple aws accounts or want it to run with a non standard profile. 

```
alias prod_ssh='AWS_PROFILE=prod sawsh'
```
