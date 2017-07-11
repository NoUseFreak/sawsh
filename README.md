# sawsh

SSH wrapper for aws to make your life easier.

## Usage

```
sawsh webserver
```

This example will query AWS for a EC2 instance containing the name `webserver`. It will prompt you with a choise when
more than one result is found.

## Install

```
git clone https://github.com/NoUseFreak/sawsh.git
cd sawsh
make
make install
```
## Configure

Setup sure your [aws-cli](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html) is setup. That user needs `ec2:Describe*` permissions.

### Suggestion

It may be useful to setup some aliases if you use multiple aws accounts or want it to run with a non standard profile. 

```
alias prod_ssh='AWS_PROFILE=prod sawsh'
```
