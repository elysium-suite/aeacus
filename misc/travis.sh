## DO NOT RUN THIS ON YOUR VM, THIS IS FOR TRAVIS ##

# Grab dependencies
echo "[+] Getting general dependencies..."
go get "github.com/urfave/cli"
go get "github.com/BurntSushi/toml/cmd/tomlv"
go get "github.com/fatih/color"

# Add convenient functions for building
echo "[+] Adding functions..."

# aeacus-build-linux-production --> build aeacus and phocus, stripped
ablp() {
  echo "Building aeacus..."
  go build -ldflags '-s -w' -o ./aeacus ./src
  echo "Linux aeacus build successful!"

  echo "Building phocus..."
  go build -ldflags '-s -w' -tags phocus -o ./phocus ./src
  echo "Linux phocus build successful!"
}

# aeacus-build-windows-production --> build aeacus and phocus, stripped (for windows)
abwp() {
  echo "Building aeacus..."
  GOOS=windows go build -ldflags '-s -w' -o ./aeacus.exe ./src
  echo "Windows aeacus build successful!"

  echo "Building phocus..."
  GOOS=windows go build -ldflags '-s -w' -tags phocus -o ./phocus.exe ./src
  echo "Windows phocus build successful!"
}

# Windows dependencies (will cause errors on Linux systems due to build constraints)
echo "[+] Getting Windows-specific dependencies..."
go get "github.com/iamacarpet/go-win64api"
go get "github.com/go-ole/go-ole"
go get "golang.org/x/sys/windows"
go get "github.com/gen2brain/beeep"
go get "github.com/go-toast/toast"
go get "github.com/tadvi/systray"
go get "github.com/judwhite/go-svc/svc"

ablp
abwp
