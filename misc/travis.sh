# Grab dependencies
echo "[+] Getting general dependencies..."
go get "github.com/urfave/cli"
go get "github.com/BurntSushi/toml/cmd/tomlv"
go get "github.com/fatih/color"

# Add convenient aliases for building
echo "[+] Adding aliases..."

# aeacus-build-linux --> build aeacus and phocus
aeacus_build_linux() {
  cd aeacus-src
  go build -o ../aeacus .
  cd ..
  cd phocus-src
  go build -o ../phocus .
  cd ..
}

# aeacus-build-linux-production --> build aeacus and phocus, stripped
aeacus_build_linux_production() {
  cd aeacus-src
  go build -ldflags '-s -w' -o ../aeacus .
  cd ..
  cd phocus-src
  go build -ldflags '-s -w' -o ../phocus .
  cd ..
}

# aeacus-build-windows --> build aeacus and phocus (for windows)
aeacus_build_windows() {
  cd aeacus-src
  GOOS=windows go build -o ../aeacus.exe .
  cd ..
  cd phocus-src
  GOOS=windows go build -o ../phocus.exe .
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
