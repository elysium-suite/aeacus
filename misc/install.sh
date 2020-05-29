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
apt install -y golang-go git

# Grab dependencies
go get "github.com/urfave/cli"
go get "github.com/BurntSushi/toml/cmd/tomlv"
go get "github.com/fatih/color"

# Windows dependencies
go get "github.com/iamacarpet/go-win64api"
go get "github.com/go-ole/go-ole"
go get "golang.org/x/sys/windows"
go get "github.com/gen2brain/beeep"
go get "github.com/go-toast/toast"
go get "github.com/tadvi/systray"
go get "github.com/judwhite/go-svc/svc"

# Add convenient aliases for building

    # builda --> build aeacus, buildp --> build phocus
    echo "alias builda=\"cd aeacus-src; go build -o ../aeacus .; cd ..\"; alias buildp=\"cd phocus-src; go build -o ../phocus .; cd ..\"" >> /etc/bash.bashrc

    # pbuilda --> builda aeacus for production, pbuildp --> build production phocus
    echo "alias pbuilda=\"cd aeacus-src; go build -ldflags '-s -w' -o ../aeacus .; cd ..\"; alias pbuildp=\"cd phocus-src; go build -ldflags '-s -w' -o ../phocus .; cd ..\"" >> /etc/bash.bashrc

# Add convenient aliases for building for Windows

    # builda --> build aeacus, buildp --> build phocus
    echo "alias wbuilda=\"cd aeacus-src; GOOS=windows go build -o ../aeacus.exe .; cd ..\"; alias wbuildp=\"cd phocus-src; GOOS=windows go build -o ../phocus.exe .; cd ..\"" >> /etc/bash.bashrc

    # pbuilda --> builda aeacus for production, pbuildp --> build production phocus
    echo "alias wpbuilda=\"cd aeacus-src; GOOS=windows go build -ldflags '-s -w' -o ../aeacus.exe .; cd ..\"; alias wpbuildp=\"cd phocus-src; GOOS=windows go build -ldflags '-s -w' -o ../phocus.exe .; cd ..\"" >> /etc/bash.bashrc

# Source aliases from /etc/bash.bashrc
source /etc/bash.bashrc
