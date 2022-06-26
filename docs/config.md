# aeacus

## Fields

This is a list of (non-check) image configuration fields for `aeacus`. For details on check configurations (ex., ensure this file has this content), see the [checks configuration](./checks.md).

**name**: Image name, primarily used to organize remote scoring.

> **Note!** This field is not mandatory if you use local scoring.

```
name = "ubuntu-18-dabbingdabbers"
```

**title**: Round title, shown in the image's scoring report and README.

```
title = "CyberPatio Practice Round 1337"
```

**os**: Name of the operating system, shown in the image's README.

```
os = "TempleOS 5.03"
```

**user**: Main user of the image. This is used when sending notifications.

> **Note!** No other user accounts will get notifications except for this user.

```
user = "sysadmin"
```

**remote**: Address of remote server for scoring. If remote scoring is enabled, and `local` is not enabled, `aeacus` will refuse to score the image unless a connection to the server can be established.

```
remote = "http://scoring.example.com"
```

**password**: Password used for encrypting remote reporting traffic. The same password must be set on the remote side.

```
password = "H4!b5at+kWls-8yh4Guq"
```

**local**: Enables local scoring. If no remote address is specified, this will automatically be set to true.

```
local = true
```

**enddate**: Defines competition end date. If the engine is run after this date, it will not score the image.

```
enddate = "2004/06/05 13:09:00 PDT"
```

**shell**: Determines if remote shell functionality is enabled. This is disabled by default. If enabled, competition organizers can interact with images from the scoring endpoint

```
shell = true
```

**version**: Version of aeacus that the configuration was made for. Used for compatibility checks, the engine will throw a warning if the binary version does not match the version specified in this field. You should set this to the version of aeacus you are using.

```
version = "X.X.X"
```

## Penalties

Assign a check a negative point value, and it will become a penalty. Example:

```
[[check]]
message = "Critical service OpenSSH stopped or removed"
points = "-5"

    [[check.passoverride]]
    type = 'ServiceUpNot'
    name = 'sshd'

    [[check.passoverride]]
    type = 'PathExistsNot'
    name = '/lib/systemd/system/sshd.service'
```

## Combining check conditions

Using multiple conditions for a check can be confusing at first, but can greatly improve the quality of your images by accounting for edge cases and abuse.

Given no conditions, a check does not pass.

**Pass** conditions act as a logical AND with other pass conditions. This means they must ALL be true for a check to pass.

**PassOverride** conditions act as a logical OR. This means that any can succeed for the check to pass.

If any **Fail** conditions succeed, the check does not pass.

So, it's like: ``((all pass checks) OR passoverride) AND fails``.

For example:

```
[[check]]

    # Pass only if both scheduled tasks are deleted
    [[check.pass]]
    type = 'ScheduledTaskExistsNot'
    name = 'Disk Cleanup'
    [[check.pass]]
    type = 'ScheduledTaskExistsNot'
    name = 'Disk Cleanup Backup'

    # OR if the user runnning those tasks is deleted
    [[check.passoverride]]
    type = 'UserExistsNot'
    name = 'CleanupBot'

    # AND the scheduled task service is running
    [[check.fail]]
    type = 'ServiceUpNot'
    name = 'Schedule'
```

The evaluation of checks goes like this:
1. Check if any Fail are true. If any Fail checks succeed, then we're done, the check doesn't pass.
2. Check if any PassOverride conditions pass. If they do, we're done, the check passes.
3. Check status of all Pass conditions. If they all succeed, the check passes, otherwise it fails.
