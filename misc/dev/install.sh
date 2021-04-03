#!/bin/sh

cat <<-EOF

 .oooo.    .ooooo.   .oooo.    .ooooo.  oooo  oooo   .oooo.o
\`P  )88b  d88' \`88b \`P  )88b  d88' \`"Y8 \`888  \`888  d88(  "8
 .oP"888  888ooo888  .oP"888  888        888   888  \`"Y88b.
d8(  888  888    .o d8(  888  888   .o8  888   888  o.  )88b
\`Y888""8o \`Y8bod8P' \`Y888""8o \`Y8bod8P'  \`V88V"V8P' 8""888P'

This script sets up the development environment on a Linux (Debian-based) box.

EOF

printf "\033[32;1m[+] Updating package lists\033[0m\n"
sudo apt-get update

printf "\n\033[32;1m[+] Installing Go\033[0m\n"
wget -O ~/go.tar.gz https://golang.org/dl/go1.16.2.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf ~/go.tar.gz

printf "\n\033[32;1m[+] Adding \`go\` binary to PATH\033[0m\n"
echo "export PATH=$PATH:/usr/local/go/bin:/$HOME/go" | sudo tee -a /etc/profile

printf "\n\033[32;1m[+] Installing Git & Make\033[0m\n"
sudo apt-get install -y git make

printf "\n\033[32;1m[+] Build dependencies installed successfully\033[0m\n"
echo "Run \`source /etc/profile\` to add \`go\` and to your PATH"
echo "Run go get -v -t -d ./... to install aeacus' dependencies"
echo "Check out the \`Makefile\` to see what targets you can build Aeacus for"
