# todo

- remote
  - status/time limit actually enforced
    - see comments in scoring.go and remote.go
- info
  - other things? esp for windows
- windows

  - improve scoring.conf example crypto (add aes-gcm, obfuscate key, etc)
  - fix windows service quit WaitGroup (phocus_windows.go)
  - ^^ THIS IS LARGELY FIXED. I count null bytes to detect unicode vs ansi. However, when the text read is only one character (for example, `b`), it will fail if unicode

- security

  - disable net/http using HTTP_PROXY environmental variable

- checks to implement

  - windows startup programs
  - windows and linux updates and auto-updating status (apt only for linux)
  - windows (Make less janky)
  - windows service-specific hardening and checks
  - windows DEP

- release

  - windows
    - Detect if firefox.exe is in x86 Program Files or just Program Files

- cool/stretch goal

  - make binary pattern in background of score report personally identifiable (like, it's their ID or something)

- hard/long term
  - verify binary
  - replace shell checks with lower-level more reliable things, winAPI, whatever
