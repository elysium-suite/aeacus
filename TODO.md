# todo

- remote
  - status/time limit actually enforced
    - see comments in scoring.go and remote.go
- gui for ID (and checks, stretch goal)
- info
  - create list of admins/users
  - other things? esp for windows
- windows
  - improve scoring.conf example crypto (add aes-gcm, obfuscate key, etc)
  - disable net/http using HTTP_PROXY environmental variable
  - fix windows service quit WaitGroup (phocus_windows.go)
  - reading TEAMID fails beacuse it's unicode by default and ioutil/program expects ANSI
  - as above, this is for FQs too. add ability to detect, or run through pruning function to remove null terms
  - ^^ THIS IS LARGELY FIXED. I count null bytes to detect unicode vs ansi. However, when the text read is only one character (for example, `b`), it will fail if unicode

 - security
    - check for ptrace attachment

- checks to implement

  - windows startup programs
  - windows features
  - windows and linux updates and auto-updating status (apt only for linux)
  - windows (Make less janky)
  - windows and linux firefox checks
  - service-specific hardening and checks
    - do we implement this for linux or just let people use filecontains?
    - do we implement sysctl checks?
  - windows DEP

- release

  - windows
    - Detect if firefox.exe is in x86 Program Files or just Program Files
    - clear run and command history
  - linux
    - clear ff history?

- hard/long term
  - verify binary
  - replace shell checks with lower-level more reliable things, winAPI, whatever
