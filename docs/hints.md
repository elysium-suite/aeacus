# Hints

Hints let you provide information on failing checks.

![Hint Example](../misc/gh/ReadMe.png)

Hints are a way to help make images more approachable.

You can add a conditional hint or a check-wide hint. A conditional hint is printed when the condition is executed and fails. Make sure you understand the check precedence; this can be trickey, as sometimes your check is NOT executed ([read about conditions](conditions.md)).

Example conditional hint:
```
[[check]]
points = 5

	[[check.pass]]
	type = "ProgramInstalledNot"
	name = "john"

	[[check.pass]]
	# This hint will NOT print unless the condition above succeeds.
	# Pass conditions are logically AND-- they all need to succeed.
	# If one fails, there's no reason to execute the other ones.
	hint = "Removing just the binary is insufficient; use a package manager to remove all of a tool's files."
	type = "PathExistsNot"
	path = "/usr/share/john"
```

Check-wide hints are at the top level and always displayed if a check fails. Example check-wide hint:

```
[[check]]
hint = "Are there any 'hacking' tools installed?"
points = 5

	[[check.pass]]
	type = "ProgramInstalledNot"
	name = "john"

	[[check.pass]]
	type = "PathExistsNot"
	path = "/usr/share/john"
```
	
You can combine check-wide and conditional hints. If the check fails, the check-wide hint is ALWAYS displayed, in addition to any conditional hints triggered.
