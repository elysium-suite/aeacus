#!/bin/sh

cat <<EOF

 .oooo.    .ooooo.   .oooo.    .ooooo.  oooo  oooo   .oooo.o
\`P  )88b  d88' \`88b \`P  )88b  d88' \`"Y8 \`888  \`888  d88(  "8
 .oP"888  888ooo888  .oP"888  888        888   888  \`"Y88b.
d8(  888  888    .o d8(  888  888   .o8  888   888  o.  )88b
\`Y888""8o \`Y8bod8P' \`Y888""8o \`Y8bod8P'  \`V88V"V8P' 8""888P'

This script sets up the development environment on a Linux (Debian-based) box.

EOF

[ "$(id -u)" = 0 ] || {
	echo "Please run this script as root!"
	exit 1
}

echo "[+] Updating package lists"
apt-get update

echo "[+] Installing Go"
wget -O ~/go.tar.gz https://golang.org/dl/go1.16.2.linux-amd64.tar.gz
tar -C /usr/local -xzf ~/go.tar.gz

echo "Adding \`go\` binary to PATH"
echo "export PATH=$PATH:/usr/local/go/bin/$HOME/go" >>/etc/profile

echo "[+] Installing Git & Make"
apt-get install -y git make

echo "[+] Installing Garble"
go get -u mvdan.cc/garble

echo "[+] Build dependencies installed successfully"
echo "Run \`source /etc/profile\` to add \`go\` and \`garble\` to your PATH"
echo "Run go get -v -t -d ./... to install aeacus' dependencies"
echo "Check out the \`Makefile\` to see what targets you can build Aeacus for"
