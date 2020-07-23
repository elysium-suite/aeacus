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
  echo "Please run this script as root!"
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

# Grab dependencies
echo "[+] Getting general dependencies..."
go get "github.com/urfave/cli"
go get "github.com/BurntSushi/toml/cmd/tomlv"
go get "github.com/fatih/color"

# Add convenient aliases for building
if ! grep -q "aeacus-build" /etc/bash.bashrc; then
  echo "[+] Adding aliases..."

  # aeacus-build-linux --> build aeacus and phocus
  echo "alias aeacus-build-linux=\"go build -o ./aeacus ./src; go build -tags phocus -o ./phocus ./src\"" >>/etc/bash.bashrc

  # aeacus-build-linux-production --> build aeacus and phocus, stripped
  echo "alias aeacus-build-linux-production=\"go build -ldflags '-s -w ' -o ./aeacus ./src; go build -ldflags '-w -s' -tags phocus -o  ./phocus ./src\"" >>/etc/bash.bashrc

  # aeacus-build-windows --> build aeacus and phocus (for windows)
  echo "alias aeacus-build-windows=\"GOOS=windows go build -o ./aeacus.exe ./src; GOOS=windows go build -tags phocus -o ./phocus.exe ./src\"" >>/etc/bash.bashrc

  # aeacus-build-windows-production --> build aeacus and phocus, stripped
  echo "alias aeacus-build-windows-production=\"GOOS=windows go build -ldflags '-s -w ' -o ./aeacus.exe ./src; GOOS=windows go build -ldflags '-w -s' -tags phocus -o ./phocus.exe ./src\"" >>/etc/bash.bashrc
fi

# Windows dependencies (will cause errors on Linux systems due to build constraints)
echo "[+] Getting Windows-specific dependencies..."
go get "github.com/iamacarpet/go-win64api"
go get "github.com/go-ole/go-ole"
go get "golang.org/x/sys/windows"
go get "github.com/gen2brain/beeep"
go get "github.com/go-toast/toast"
go get "github.com/tadvi/systray"
go get "golang.org/x/text/unicode"
go get "github.com/judwhite/go-svc/svc"

source /etc/bash.bashrc
