.DEFAULT_GOAL := linux-secure

linux:
	go build -o ./aeacus ./src; go build -tags phocus -o ./phocus ./src

linux-secure:
	go build -ldflags '-s -w ' -o ./aeacus ./src; go build -ldflags '-w -s' -tags phocus -o  ./phocus ./src

windows:
	GOOS=windows go build -o ./aeacus.exe ./src; GOOS=windows go build -tags phocus -o ./phocus.exe ./src

windows-secure:
	GOOS=windows go build -ldflags '-s -w ' -o ./aeacus.exe ./src; GOOS=windows go build -ldflags '-w -s' -tags phocus -o ./phocus.exe ./src
