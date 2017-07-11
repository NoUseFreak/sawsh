# sawsh

SSH wrapper for aws to make your life easier.

## Usage

```
sawsh webserver
```

This example will query AWS for a EC2 instance containing the name `webserver`. It will prompt you with a choise when
more than one result is found.

```
$ sawsh webserver
listing instances with tag core in: us-east-1
+---+--------------------+-------------+
|   |       NAME         |      IP     |
+---+--------------------+-------------+
| 0 | prod-webserver-1   | 10.1.1.10   |
| 1 | prod-webserver-2   | 10.1.2.10   |
| 2 | prod-webserver-3   | 10.1.3.10   |
| 3 | prod-webserver-4   | 10.1.1.11   |
| 4 | prod-webserver-5   | 10.1.2.11   |
| 5 | prod-webserver-6   | 10.1.3.11   |
+---+--------------------+-------------+
Pick a number: 1
Connecting to 10.1.2.10 ...
```

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
