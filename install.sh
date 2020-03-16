############################################################
cat << EOF

 .oooo.    .ooooo.   .oooo.    .ooooo.  oooo  oooo   .oooo.o
\`P  )88b  d88' \`88b \`P  )88b  d88' \`"Y8 \`888  \`888  d88(  "8
 .oP"888  888ooo888  .oP"888  888        888   888  \`"Y88b.
d8(  888  888    .o d8(  888  888   .o8  888   888  o.  )88b
\`Y888""8o \`Y8bod8P' \`Y888""8o \`Y8bod8P'  \`V88V"V8P' 8""888P'

EOF
############################################################

apt update
apt install -y software-properties-common
add-apt-repository ppa:longsleep/golang-backports
apt update
apt install -y golang-go git
go get "github.com/urfave/cli"
go get "github.com/BurntSushi/toml/cmd/tomlv"
go get "github.com/fatih/color"
echo "alias builda=\"cd aeacus-src; go build -o ../aeacus .; cd ..\"; alias buildp=\"cd phocus-src; go build -o ../phocus .; cd ..\""
