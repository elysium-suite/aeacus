# aeacus

## Adding Crypto

The public releases of `aeacus` ship with weak crypto (cryptographic security), which means that the encryption and/or encoding of scoring data files is not very "secure".

You can compile it yourself to generate random keys (`make release`). This means the public release decrypt function will not work.

If security of the configuration is very important to you, or you feel the competition integrity is at risk, (e.g., you're running a competition with prizes, or running a practice session for beginner reverse engineers), you should compile the binary for yourself after adding stronger crypto operations.

This is not too hard. Anything you want to add is good! More XOR, AES (be careful to implement this one correctly), mixing bytes up... As long as the encrypt and decrypt functions work, nothing should break. The functions you would want to change are all in `crypto.go`.

Once you implement your functions, make sure they work. You can run built-in tests with:

```bash
CGO_ENABLED=0 go test -v
```

This model of engine can never be 100% secure (see [security](security.md)), but you can get pretty ok security.
