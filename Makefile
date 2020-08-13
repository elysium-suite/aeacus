aeacus-build-linux:
	go build -o ./aeacus ./src; go build -tags phocus -o ./phocus ./src

aeacus-build-linux-production:
	go build -ldflags '-s -w ' -o ./aeacus ./src; go build -ldflags '-w -s' -tags phocus -o  ./phocus ./src

aeacus-build-windows:
	GOOS=windows go build -o ./aeacus.exe ./src; GOOS=windows go build -tags phocus -o ./phocus.exe ./src

aeacus-build-windows-production:
	GOOS=windows go build -ldflags '-s -w ' -o ./aeacus.exe ./src; GOOS=windows go build -ldflags '-w -s' -tags phocus -o ./phocus.exe ./src
