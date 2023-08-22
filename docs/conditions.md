# Check conditions and precedence

Using multiple conditions for a check can be confusing at first, but can greatly improve the quality of your images by accounting for edge cases and abuse.

Given no conditions, a check does not pass.

If any **Fail** conditions succeed, the check does not pass.

**PassOverride** conditions act as a logical OR. This means that any can succeed for the check to pass.

**Pass** conditions act as a logical AND with other pass conditions. This means they must ALL be true for a check to pass.

If the outcome of a check is decided, aeacus will NOT execute the remaining conditions (it will "short circuit"). For example, if a PassOverride succeeds, any Pass conditions are NOT executed.

So, it's like this: `check_passes = (NOT fails) AND (passoverride OR (AND of all pass checks))`.

For example:

```
[[check]]

    # Ensure the scheduled task service is running AND
    [[check.fail]]
    type = 'ServiceUpNot'
    name = 'Schedule'

    # Pass if the user runnning those tasks is deleted
    [[check.passoverride]]
    type = 'UserExistsNot'
    name = 'CleanupBot'
    
    # OR pass if both scheduled tasks are deleted
    [[check.pass]]
    type = 'ScheduledTaskExistsNot'
    name = 'Disk Cleanup'
    [[check.pass]]
    type = 'ScheduledTaskExistsNot'
    name = 'Disk Cleanup Backup'

```

The evaluation of checks goes like this:
1. Check if any Fail are true. If any Fail checks succeed, then we're done, the check doesn't pass.
2. Check if any PassOverride conditions pass. If they do, we're done, the check passes.
3. Check status of all Pass conditions. If they all succeed, the check passes, otherwise it fails.
