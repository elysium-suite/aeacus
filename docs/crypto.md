# aeacus

## Adding Crypto

The public releases of `aeacus` ship with very weak crypto. You should compile the binary for yourself after adding stronger crypto. This is not too hard, and most of the work is done for you. See the example in `examples/example-crypto.go`.

Essentially, you need to provide the following functions:
- `func writeCryptoConfig(mc *metaConfig) string {` that reads the scoring config, encrypts it, and returns the encrypted string.
- `func readCryptoConfig(mc *metaConfig) string {` that reads the scoring config, decrypts it, and returns the plaintext string.
- `func encryptString(password string, plaintext string) string {` that encrypts a given plaintext with the password and returns a string.
- `func decryptString(password string, plaintext string) string {` that decrypts a given ciphertext with the password and returns the plaintext string.

The reason this is done is that, due to the design of the scoring engine, the easiest crypto solution is to use a symmetric, hard-coded key. You could try to do something asymmetric or wild, but practically to be able to run without a server and with a preshared config, all the data needs to be already within the engine. This means that it's very easy to reverse-engineer, especially (!) when it's open source. Our only methods, then, of protecting against decryption, is to make RE as difficult as possible by doing weird things with the hardcoded key and not distributing the source (by having people add their own).
