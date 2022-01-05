# aeacus [![Go Report Card](https://goreportcard.com/badge/github.com/elysium-suite/aeacus)](https://goreportcard.com/report/github.com/elysium-suite/aeacus)

<img align="right" width="200" src="assets/img/logo.png"/>

`aeacus` is a vulnerability scoring engine for Windows and Linux, with an emphasis on simplicity.

## V2

`aeacus` has recently been updated to version 2.0.0! To view the breaking changes, refer to [./docs/v2.md](./docs/v2.md).

## Installation

0. **Extract the release** into `/opt/aeacus` (Linux) or `C:\aeacus\` (Windows).

	> Try compiling it yourself! Or, you can [download the releases here](https://github.com/elysium-suite/aeacus/releases).

1. **Set up the environment.**

	- Put your **config** in `/opt/aeacus/scoring.conf` or`C:\aeacus\scoring.conf`.

		- _Don't have a config? See the example below._

	- Put your **README data** in `ReadMe.conf`.

2. **Check that your config is valid.**

```
./aeacus --verbose check
```

> Check out what you can do with `aeacus` with `./aeacus --help`!

3. **Score the image with the current config to verify your checks work as expected.**

```
./aeacus --verbose score
```

> The TeamID is read from `/opt/aeacus/TeamID.txt` or `C:\aeacus\TeamID.txt`.

4. **Prepare the image for release.**

> **WARNING**: This will remove `scoring.conf`. Back it up somewhere if you want to save it! It will also remove the `aeacus` executable and other sensitive files.

```
./aeacus --verbose release
```

## Screenshots

### Scoring Report:

![Scoring Report](./misc/gh/ScoringReport.png)

### ReadMe:

![ReadMe](./misc/gh/ReadMe.png)

## Features

-   Robust yet simple vulnerability scorer
-   Image preparation (cleanup, README, etc)
-   Remote score reporting

> Note: `aeacus` ships with weak crypto on purpose. You should implement your own crypto functions if you want to make it harder to crack. See [Adding Crypto](/docs/crypto.md) for more information.

## Compiling

Only Linux development environments are officially supported. Ubuntu virtual machines work great.

Make sure you have a recent version of `go` installed, as well as `git` and `make`. If you want to compile Windows and Linux, install all dependencies using `go get -v -d -t ./...`. Then to compile, use `go build`, OR make:

- Building for `Linux`: `make lin`
- Building for `Windows`: `make win`

### Development

If you're developing for `aeacus`, compile with these commands to leave debug symbols in the binaries:

- Building for `Linux`: `make lin-dev`
- Building for `Windows`: `make win-dev`

### Release Archives

You can build release archives (e.g., `aeacus-linux.zip`). These will have auto-generated `crypto.go` files.

- Building both platforms: `make release`

## Documentation

All checks (with examples and notes) [are documented here](docs/checks.md).

Other documentation:
- [Scoring Configuration](docs/config.md)
- [Crypto](docs/crypto.md)
- [Security Model](docs/security.md)
- [Windows Security Policy](docs/securitypolicy.md)

## Remote Endpoint

Set the `remote` field in the configuration, and your image will use remote scoring. If you want remote scoring, you will need to host a remote scoring endpoint. The authors of this project recommend using [sarpedon](https://github.com/elysium-suite/sarpedon). See [this example remote configuration](docs/examples/remote.conf).

## Configuration

The configuration is written in TOML. Here is a minimal example:

```toml
name = "ubuntu-18-supercool" # Image name
title = "CoolCyberStuff Practice Round" # Round title
os = "Ubuntu 18.04" # OS, used for README
user = "coolUser" # Main user for the image

# Set the aeacus version of this scoring file. Set this to the version
# of aeacus you are using. This is used to make sure your configuration,
# if re-used, is compatible with the version of aeacus being used.
#
# You can print your version of aeacus with ./aeacus version.
version = "2.0.0"

[[check]]
message = "Removed insecure sudoers rule"
points = 10

	[[check.pass]]
	type = "FileContainsNot"
	path = "/etc/sudoers"
	value = "NOPASSWD"

[[check]]
# If no message is specified, one is auto-generated
points = 20

	[[check.pass]]
	type = "FileExistsNot"
	path = "/usr/bin/ufw-backdoor"

	[[check.pass]]     # You can code multiple pass conditions, but
	type = "Command"   # they must ALL succeed for the check to pass!
	cmd  = "ufw status"

[[check]]
message = "Malicious user 'user' can't read /etc/shadow"
# If no points are specified, they are auto-calculated out of 100.

	[[check.pass]]
	type = "CommandNot"
	cmd  = "sudo -u user cat /etc/shadow"

	[[check.pass]]  		# "pass" conditions are logically AND with other pass
	type = "FileExists"		# conditions. This means they all must pass for a check
	path = "/etc/shadow"	# to be considered successful.

	[[check.passoverride]]  # If you a check to succeed if just one condition
	type = "UserExistsNot"  # passes, regardless of other pass checks, use
	user = "user"           # an override pass (passoverride). This is a logical OR.
							# passoverride is overridden by fail conditions.

	[[check.fail]]          # If any fail conditions succeed, the entire check will fail.
	type = "FileExistsNot"
	path = "/etc/shadow"

[[check]]
message = "Administrator has been removed"
points = -5 # This check is now a penalty, because it has negative points

	[[check.pass]]
	type = "UserExistsNot"
	user = "coolAdmin"

```

See more in-depth examples, including remote reporting, [here](https://github.com/elysium-suite/aeacus/tree/master/docs/examples).

## ReadMe Configuration

Put your README in `ReadMe.conf`. Here's a commented template:

```html
<!-- Put your comments/additions to the normal ReadMe here! -->
<p>
	Uncomplicated Firewall (UFW) is the only company approved Firewall for use
	on Linux machines at this time.
</p>

<!-- You can add as many <p></p> notes as you want! This HTML is simply imported into the existing ReadMe template. -->
<p>
	Congratulations! You just recruited a promising new team member. Create a
	new Standard user account named "bobbington" with a temporary password of
	your choosing.
</p>

<!-- Put your critical services here! -->
<p><b>Critical Services:</b></p>
<ul>
	<li>OpenSSH Server (sshd)</li>
	<li>Other cool service</li>
</ul>

<!-- Put your users here! -->
<h2>Authorized Administrators and Users</h2>

<pre>
<b>Authorized Administrators:</b>
coolUser (you)
	password: coolPassword
bob
	password: bob

<b>Authorized Users:</b>
coolFriend
awesomeUser
radUser
coolGuy
niceUser
</pre>
```

## Information Gathering

The `aeacus` binary supports gathering information (on **Windows** only) in cases where it's tough to gather what the scoring system can see.

Print information with `./aeacus info type` where `type` is one the following:

### Windows

-   `packages` (shows installed programs)
-   `users` (shows local users)
-   `admins` (shows local administrator users)

## Tips and Tricks

-   Easily change the branding by replacing `assets/img/logo.png`.
-   Test your scoring configuration in a loop:
``` bash
while true; do ./aeacus -v; sleep 20; done
```

## Contributing and Disclaimer

A huge thanks to the project contributors for help adding code and features, and to many others for help with feedback, usability, and finding bugs!

If you have anything you would like to add or fix, please make a pull request! No improvement or fix is too small, and help is always appreciated.

Thanks to UTSA CIAS and the CyberPatriot program for putting together such a cool competition, and for the inspiration to make this project.

This project is in no way affiliated with or endorsed by the Air Force Association, University of Texas San Antonio, or the CyberPatriot program.
