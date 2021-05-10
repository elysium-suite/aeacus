# aeacus

## Fields

This is a list of configuration fields that are required when creating `aeacus`. Well, most of them are required.

**name**: Image name, primarily used to organize remote scoring.

> **Note!** This field is not mandatory if you use local scoring.

```
name = "ubuntu-18-dabbingdabbers"
```

**title**: Round title, as seen in the scoring report and README.

```
title = "CyberPatio Practice Round 18"
```

**os**: Name of the operating system, as seen in the README.

```
os = "TempleOS 5.03"
```

**user**: Main user of the image.

```
user = "sysadmin"
```

**remote**: Address of remote server for scoring. If remote scoring is enabled, aeacus will refuse to score the image unless a connection to the server can be established.

```
remote = "8.8.8.8"
```

**password**: Password used for encrypting remote reporting traffic.

```
password = "H4!b5at+kWls-8yh4Guq"
```

**local**: Disables remote scoring. If no remote address is specified, this will automatically be set to true.

```
local = true
```

**enddate**: Defines self-destruct date. If the engine is run after this date, the image will self destruct. Formatted as YEAR/MO/DA HR:MN:SC ZONE

```
enddate = "2004/06/05 13:09:00 PDT"
```

**nodestroy**: Governs self-destruct behavior. If this is set to true, only the aeacus folder will be deleted, leaving the rest of the image intact.

```
nodestroy = true
```

**disableshell**: Enables remote shell functionality. If set to true, aeacus will not attempt to connect to the remote shell.

```
disableshell = false
```

**version**: Version of aeacus that the configuration was built for. Primarily used for compatibility checks, the engine will throw a warning if the binary version does not match the version specified in this field. In the future, this may also be used for backwards compatibility functionality.

```
version = "1.6.0"
```
