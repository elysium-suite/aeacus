# aeacus

This is a client-side scoring system meant to imitate the functionality of UTSA's CIAS CyberPatriot Scoring System (CSS) with an emphasis on simplicity. Named after the Greek myth of King Aeacus, a judge of the dead.

## Installation

0. Download the most recent zip from releases into /opt/ on your vulnerable virtual machine.
```
cd /opt && git clone https://github.com/sourque/aeacus/releases...
```
1. Write your config in `/opt/aeacus/scoring.conf`.
> Don't have a config? See the example at the bottom of this README.

2. Check that your config is valid.
```
aeacus --verbose check
```
3. Prepare the image for release.
```
aeacus --verbose release
```

## Screenshot

![Scoring Report](assets/img/scoring_report.png)

## Features

- In-depth and simple vulnerability scorer
- Image deployment (cleanup, README, etc)
- Remote score reporting through a REST-like API

## Checks

This is a list of vulnerability checks that can be used in the configuration for aeacus.

> __Note!__ Each of these check types can be used for either `Pass` or `Fail` conditions, and there can be multiple `Pass` or `Fail` conditions per check.

__Command__: pass if command succeeds (exit code `0`)
```
type="Command"
arg1="grep 'pam_history.so' /etc/pam.d/common-password"
```

> __Note!__ Each of the commands here can check for the opposite by appending "Not" to the check type. For example, `CommandNot` to pass if a command does not return exit code `0`.

__FileExists__: pass if specified file exists
```
type="FileExists"
arg1="/etc/passwd"
```

__FileContains__: pass if file contains string (regex enabled)
```
type="FileContains"
arg1="/home/ballen/Desktop/Forensic Question 1.txt"
arg2="ANSWER:\sBarry[a-zA-Z]+Allen"
```

__FileEquals__: pass if file equals sha1 hash
```
type="FileEquals"
arg1="/etc/sysctl.conf"
arg2="403926033d001b5279df37cbbe5287b7c7c267fa"
```

__PackageInstalled__: pass if package is installed
```
type="PackageInstalled"
arg1="tcpd"
```

__UserExists__: pass if user exists on system
```
type="UserExists"
arg1="ballen"
```

> __Note!__ If a check has negative points assigned to it, it automatically becomes a penalty.

## Configuration

The configuration is written in TOML. See the below example for all possible features:

```
example config
```

## Contributing and Disclaimer

If you have anything you would like to add or fix, please make a pull request! No improvement or fix is too small, and help is always appreciated.

This project is in no way affiliated with or endorsed by the Air Force Association, University of Texas San Antonio, or the CyberPatriot program.
