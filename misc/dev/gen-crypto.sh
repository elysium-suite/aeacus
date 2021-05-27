#!/bin/sh

rand() { openssl rand -hex 64; }
replace() { sed -i "s/$1/$2/g" cmd/crypto.go.tmpl; }

hashOne=$(rand)
hashTwo=$(rand)
byteKey=$(rand | sed 's/\(..\)/0x\1, /g')

cp cmd/crypto.go.tmpl cmd/crypto.go.tmpl.bak

replace "HASH_ONE" "$hashOne"
replace "HASH_TWO" "$hashTwo"
replace "BYTE_KEY" "$byteKey"

cp cmd/crypto.go.tmpl cmd/crypto.go
cp cmd/crypto.go.tmpl.bak cmd/crypto.go.tmpl
rm cmd/crypto.go.tmpl.bak

echo "generated crypto.go"
