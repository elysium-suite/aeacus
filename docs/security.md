# Security

Engines that work like `aeacus` can never be "secure." We can only make it more difficult to crack.

As long as the configuration is being loaded onto a virtual machine that is controlled by a competitor, there is no way to make it impossible to reveal what resources are being collected (files, command output, registry keys...). This is because they have control of the disk and CPU.

Due to this fact, we focus on obfuscating and encrypting the configuration such that it would take at least a fair bit of reversing expertise and time in order to crack. The primary target audience for this style of competition will likely not be able or willing to do that.

If you need a perfectly secure engine, there is no such thing. It is possible to do slightly better than our current architecture by either:

1. Only allowing remote scoring, and performing all checks on the remote host. (e.g., send the entire /etc/ssh/sshd_config file to the scoring server, who then checks it for a desired string.) That way, the competitor would know what resources are being collected, but not what the checks are looking for. However, this costs a lot more bandwidth and server CPU time.
2. Use "zero-knowledge" crypto algorithms. A simple version of this might be checking each substring of a file for a hash. Downsides to this approach include a much larger demand on client CPUs, and it's not perfect, since competitors can attempt to crack that hash without the scoring server in the loop.

If you need an almost perfectly secure engine, you will need to write your own, and it will need to have a couple things:

- Cloud-based VMs on controlled and monitored in-person thin-client kiosks with no external internet access.
- Scoring engine that works at the VM infrastructure level, and reads in competitor VM disk files to score images.

And even then, it's definitely breakable given some time.

In any case, `aeacus` crypto that you compile yourself, or with a few tweaks, is probably suitable for your use case.
