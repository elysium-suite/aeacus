# Security

Engines that work like `aeacus` can never be "secure." We can only make it more difficult to crack.

As long as the configuration is being loaded onto a virtual machine that is controlled by a competitor, there is no way to make it impossible to reveal what checks are being run. This is because they have control of the disk and CPU.

Due to this fact, we focus on obfuscating and encrypting the configuration such that it would take at least a fair bit of reversing expertise and time in order to crack. The primary target audience for this style of competition will likely not be able or willing to do that.

If you need a perfectly secure engine, there is no such thing. If you need an almost perfectly secure engine, you will need to write your own, and it will need to have two things:

- Cloud-based VMs on controlled and monitored in-person thin-client kiosks with no external internet access.
- Scoring engine that works at the VM infrastructure level, and reads in competitor VM disk files to score images.

And even then, it's definitely breakable given some time. In any case, `aeacus` crypto with a few tweaks is probably suitable for your use case.
