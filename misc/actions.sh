## DO NOT RUN THIS ON YOUR VM, THIS IS FOR ACTIONS ##

# Add convenient functions for building
echo "[+] Adding functions..."

# aeacus-build-linux-production --> build aeacus and phocus, stripped
ablp() {
  echo "Building aeacus..."
  go build -ldflags '-s -w' -o ./aeacus .
  echo "Linux aeacus build successful!"

  echo "Building phocus..."
  go build -ldflags '-s -w' -tags phocus -o ./phocus .
  echo "Linux phocus build successful!"
}

# aeacus-build-windows-production --> build aeacus and phocus, stripped (for windows)
abwp() {
  echo "Building aeacus..."
  GOOS=windows go build -ldflags '-s -w' -o ./aeacus.exe .
  echo "Windows aeacus build successful!"

  echo "Building phocus..."
  GOOS=windows go build -ldflags '-s -w' -tags phocus -o ./phocus.exe .
  echo "Windows phocus build successful!"
}

ablp
abwp
