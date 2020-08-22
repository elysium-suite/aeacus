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
source /etc/profile

# Install git (for go get)
echo "[+] Installing git..."
apt install -y git

# Add convenient aliases for building
if ! grep -q "aeacus-build" /etc/profile; then
    echo "[+] Adding aliases..."

    # aeacus-build-linux --> build aeacus and phocus, stripped
    echo "alias aeacus-build-linux=\"go build -ldflags '-s -w '; go build -ldflags '-w -s' -tags phocus -o  ./phocus\"" >> /etc/profile

    # aeacus-build-windows --> build aeacus and phocus, stripped
    echo "alias aeacus-build-windows=\"GOOS=windows go build -ldflags '-s -w '; GOOS=windows go build -ldflags '-w -s' -tags phocus -o ./phocus.exe\"" >> /etc/profile
fi

# Windows dependencies (will cause errors on Linux systems due to build constraints)

echo "[+] Make sure to start a new session or source /etc/profile!"
echo "[+] Done!"
