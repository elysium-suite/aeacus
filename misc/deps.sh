echo "[+] Getting general dependencies..."
go get "github.com/urfave/cli"
go get "github.com/BurntSushi/toml/cmd/tomlv"
go get "github.com/fatih/color"

echo "[+] Getting Windows-specific dependencies..."
go get "github.com/iamacarpet/go-win64api"
go get "github.com/go-ole/go-ole"
go get "golang.org/x/sys/windows"
go get "github.com/gen2brain/beeep"
go get "github.com/go-toast/toast"
go get "github.com/tadvi/systray"
go get "github.com/judwhite/go-svc/svc"

echo "[+] Running Unit Tests"
go test ./src
