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

**DisableRemoteEncryption**: Disables encryption of remote reporting traffic. This is not recommended, but can be useful for debugging or if you are using a custom remote endpoint that does not support encryption.

```
DisableRemoteEncryption = true
```

**local**: Enables local scoring. If no remote address is specified, this will automatically be set to true.

```
local = true
```

**enddate**: Defines competition end date. If the engine is run after this date, it will not score the image.

```
enddate = "2004/06/05 13:09:00 PDT"
```

**shell**: (Warning: the canonical remote endpoint (sarpedon) does not support this feature). Determines if remote shell functionality is enabled. This is disabled by default. If enabled, competition organizers can interact with images from the scoring endpoint

```
shell = false
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


