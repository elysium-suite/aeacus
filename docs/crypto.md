# aeacus

## Adding Crypto

The public releases of `aeacus` ship with very weak crypto. You should compile the binary for yourself after adding stronger crypto. This is not too hard, and almost all of the work is done for you.

**See the example in `aeacus-src/crypto.go`.** The project should compile as is, but again, you should read and change the file at least a little bit, then compile it for yourself.

You must specify the `remoteBackupKey` variable in your `crypto.go` and implement `writeCryptoConfig(mc *metaConfig)` and `readCryptoConfig(mc *metaConfig)`.
