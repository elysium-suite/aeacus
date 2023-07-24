# Remote score reporting

`aeacus` can report scores to a remote endpoint if the `remote` field is set in the configuration.

After every scoring cycle, `aeacus` sends an update with the following structure to the `/update` endpoint on the remote URL:

```
IMAGE_POINTS
TOTAL_VULNS
Penalty message - 10 pts
Vuln message - 5 pts
Another vuln message - 3 pts
```

So a real report may look like:

```
10
22
Firewall has been activated - 5 pts
```

Each newline in the text above actually represents a delimiter, which is a sequence of two bytes (0xff followed by 0xde). This was randomly chosen and serves to separate information from one another in the config update more reliably than a newline.

The report is encrypted with the configuration password, and hex encoded, before being sent to the remote.

The function `genUpdate()` creates the report while `reportScore()` sends it.
