#!/bin/sh
set -e

rand() { xxd -l 64 -c 64 -p /dev/urandom; }
replace() { sed -i "s/$1/$2/g" crypto.go; }

hashOne=$(rand)
hashTwo=$(rand)
byteKey=$(rand | sed 's/\(..\)/0x\1, /g')

if [ -f "crypto.go.bak" ]; then
	mv crypto.go.bak crypto.go
fi

cp crypto.go crypto.go.bak

replace "HASH_ONE" "$hashOne"
replace "HASH_TWO" "$hashTwo"
replace "0x01" "`echo $byteKey | head -c 382`"

echo "generated crypto.go"
