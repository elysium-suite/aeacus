############################################################
cat << EOF

 .oooo.    .ooooo.   .oooo.    .ooooo.  oooo  oooo   .oooo.o
\`P  )88b  d88' \`88b \`P  )88b  d88' \`"Y8 \`888  \`888  d88(  "8
 .oP"888  888ooo888  .oP"888  888        888   888  \`"Y88b.
d8(  888  888    .o d8(  888  888   .o8  888   888  o.  )88b
\`Y888""8o \`Y8bod8P' \`Y888""8o \`Y8bod8P'  \`V88V"V8P' 8""888P'

EOF
############################################################

# This script sets up the development environment on a Linux (Debian-based) box.

# Update package list
apt update

# Install golang and git (for go get)
echo "[+] Installing go and git..."
apt install -y golang-go git

# Grab dependencies
echo "[+] Getting general dependencies..."
go get "github.com/urfave/cli"
go get "github.com/BurntSushi/toml/cmd/tomlv"
go get "github.com/fatih/color"
go get "github.com/gen2brain/beeep"

# Add convenient aliases for building
if ! grep -q "aeacus-build" /etc/bash.bashrc; then
    echo "[+] Adding aliases..."

    # aeacus-build-linux --> build aeacus and phocus
    echo "alias aeacus-build-linux=\"cd aeacus-src; go build -o ../aeacus .; cd ..; cd phocus-src; go build -o ../phocus .; cd ..\"" >> /etc/bash.bashrc

    # aeacus-build-linux-production --> build aeacus and phocus, stripped
    echo "alias aeacus-build-linux-production=\"cd aeacus-src; go build -ldflags '-s -w' -o ../aeacus .; cd ..; cd phocus-src; go build -ldflags '-s -w' -o ../phocus .; cd ..\"" >> /etc/bash.bashrc

    # aeacus-build-windows --> build aeacus and phocus (for windows)
    echo "alias aeacus-build-windows=\"cd aeacus-src; GOOS=windows go build -o ../aeacus.exe .; cd ..; cd phocus-src; GOOS=windows go build -o ../phocus.exe .; cd ..\"" >> /etc/bash.bashrc

    # aeacus-build-windows-production --> build aeacus and phocus, stripped
    echo "alias aeacus-build-windows-production=\"cd aeacus-src; GOOS=windows go build -ldflags '-s -w' -o ../aeacus.exe .; cd ..; cd phocus-src; GOOS=windows go build -ldflags '-s -w' -o ../phocus.exe .; cd ..\"" >> /etc/bash.bashrc

fi

# Source aliases from /etc/bash.bashrc
source /etc/bash.bashrc

# Windows dependencies (will cause errors on Linux systems due to build constraints)
echo "[+] Getting Windows-specific dependencies..."
go get "github.com/iamacarpet/go-win64api"
go get "github.com/go-ole/go-ole"
go get "golang.org/x/sys/windows"
go get "github.com/gen2brain/beeep"
go get "github.com/go-toast/toast"
go get "github.com/tadvi/systray"
go get "github.com/judwhite/go-svc/svc"
