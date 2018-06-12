# sawsh
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FNoUseFreak%2Fsawsh.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2FNoUseFreak%2Fsawsh?ref=badge_shield)
[![Build Status](https://travis-ci.org/NoUseFreak/sawsh.svg?branch=master)](https://travis-ci.org/NoUseFreak/sawsh)

SSH wrapper for aws to make your life easier.

## Features

 - Lookup AWS EC2 instances by name.
 - Lookup ip by `ip-xxx.xxx.xxx.xxx` format.
 - Transparant connect to ip

## Usage

```sh
$ sawsh webserver
```

This example will query AWS for a EC2 instance containing the name `webserver`. It will prompt you with a choise when
more than one result is found.

```sh
$ sawsh webserver
+---+--------------------+-------------+-------------------------------+
|   |        NAME        |      IP     |          LAUNCHTIME           |
+---+--------------------+-------------+-------------------------------+
| 0 | prod-webserver-1   | 10.1.1.10   | 2018-02-01 21:13:44 +0000 UTC |
| 1 | prod-webserver-2   | 10.1.2.10   | 2018-03-15 18:57:02 +0000 UTC |
| 2 | prod-webserver-3   | 10.1.3.10   | 2018-04-19 18:04:07 +0000 UTC |
| 3 | prod-webserver-4   | 10.1.1.11   | 2018-02-15 12:36:45 +0000 UTC |
| 4 | prod-webserver-5   | 10.1.2.11   | 2018-06-07 15:54:00 +0000 UTC |
| 5 | prod-webserver-6   | 10.1.3.11   | 2018-06-07 15:54:00 +0000 UTC |
+---+--------------------+-------------|-------------------------------+
Pick a number: 1
Connecting to 10.1.2.10 ...
```

```
$ sawsh -h
NAME:
   sawsh - Query and connect to ec2 instances

USAGE:
   sawsh [global options] command [command options] [filter]

COMMANDS:
     connect  Search and connect to an instance
     exec     Search and execute a command on multiple servers
     list     Render a list of instances
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --aws-region value  (default: "us-east-1")
   --help, -h          show help

COPYRIGHT:
   (c) Dries De Peuter <dries@depeuter.io>

```

## Install

### Official release

Download the latest [release](https://github.com/NoUseFreak/sawsh/releases).

```bash
curl -sL http://bit.ly/get-sawsh | bash
```

### Build from source

```sh
$ git clone https://github.com/NoUseFreak/sawsh.git
$ cd sawsh
$ make
$ make install
```

### Upgrade

To upgrade to the latest repeat the install step.

## Configure

Setup sure your [aws-cli](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html) is setup. That user needs `ec2:Describe*` permissions.

### Suggestion

It may be useful to setup some aliases if you use multiple aws accounts or want it to run with a non standard profile. 

```sh
$ alias prod_ssh='AWS_PROFILE=prod sawsh'
```


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FNoUseFreak%2Fsawsh.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2FNoUseFreak%2Fsawsh?ref=badge_large)
