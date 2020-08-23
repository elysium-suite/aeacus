############################################################
cat <<EOF

 .oooo.    .ooooo.   .oooo.    .ooooo.  oooo  oooo   .oooo.o
\`P  )88b  d88' \`88b \`P  )88b  d88' \`"Y8 \`888  \`888  d88(  "8
 .oP"888  888ooo888  .oP"888  888        888   888  \`"Y88b.
d8(  888  888    .o d8(  888  888   .o8  888   888  o.  )88b
\`Y888""8o \`Y8bod8P' \`Y888""8o \`Y8bod8P'  \`V88V"V8P' 8""888P'

EOF
############################################################

# This script sets up the development environment on a Linux (Debian-based) box.

# Force script to be run as root
if [ "$EUID" -ne 0 ]; then
  echo "Please run this script as root! It's very short-- please feel free to audit its source code."
  exit 1
fi

# Update package list
apt update

# Install golang
echo "[+] Installing golang..."
wget -O ~/go1.14.5.linux-amd64.tar.gz https://golang.org/dl/go1.14.5.linux-amd64.tar.gz
tar -C /usr/local -xzf ~/go1.14.5.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >>/etc/profile

# Install git (for go get)
echo "[+] Installing git..."
apt install -y git

# Finalize
echo "[+] Dependencies installed successfully!"
echo "Run \`source /etc/profile\` to add \`go\` to your PATH"
echo "Check out the \`Makefile\` to see what targets you can build Aeacus for!"
