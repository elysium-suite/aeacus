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
  - binary reg checks

- security

  - rsa pub/privkey infra for encrypting scoring config !!! (thanks alvin)
  - disable net/http using HTTP_PROXY environmental variable
  - obfuscate nonencrypted config args to obfuscate types of call
    - right now they're empty
  - add fake reg/file retrives to obfuscate real calls on windwos
  - add fake file retrieves/read to obfuscate calls on linux
  - anti-tracing and anti-debugging
    - delete system if tracing detected
    - use syscall PTRACEME and see if it errors out (or similar)
    - check if any ebpf blobs are loaded into kernel
    - refuse to run if not signed (? how to implement)
 - verified vulns on other side (have vuln list on server as well)

- QoL
    - spellcheck/typo alert in config
    - if arg number is wrong, alert them (ex. no arg2 when its required)
    - TESTS!!!
        - TESTS!!!
        - TESTS!!!
        - TESTS!!!
        - TESTS!!!
        - TESTS!!!
        - TESTS!!!
        - TESTS!!!

- checks to implement

  - windows startup programs
  - windows and linux updates
  - windows (Make less janky)
  - windows service-specific hardening and checks
  - windows DEP

- release

  - windows
    - Detect if firefox.exe is in x86 Program Files or just Program Files
    - clear regedit opening

- cool/stretch goal

  - make binary pattern in background of score report personally identifiable (like, it's their ID or something)

- hard/long term
  - verify binary
  - replace shell checks with lower-level more reliable things, winAPI, whatever
