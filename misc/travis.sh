## DO NOT RUN THIS ON YOUR VM, THIS IS FOR TRAVIS ##

# Grab dependencies
echo "[+] Getting general dependencies..."
go get "github.com/urfave/cli"
go get "github.com/BurntSushi/toml/cmd/tomlv"
go get "github.com/fatih/color"

# Add convenient aliases for building
echo "[+] Adding aliases..."

# aeacus-build-linux-production --> build aeacus and phocus, stripped
aeacus_build_linux_production() {
  cd aeacus-src
  go build -ldflags '-s -w' -o ../aeacus .
  cd ..
  cd phocus-src
  go build -ldflags '-s -w' -o ../phocus .
  cd ..
}

# aeacus-build-windows-production --> build aeacus and phocus, stripped (for windows)
aeacus_build_windows_production() {
  cd aeacus-src
  GOOS=windows go build -ldflags '-s -w' -o ../aeacus.exe .
  cd ..
  cd phocus-src
  GOOS=windows go build -ldflags '-s -w' -o ../phocus.exe .
  cd ..
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

aeacus_build_linux_production
aeacus_build_windows_production
