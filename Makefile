.DEFAULT_GOAL := all

all:
	GOOS=windows garble build -ldflags '-s -w' -tags phocus -o ./phocus.exe . && GOOS=windows garble build -ldflags '-s -w' -o ./aeacus.exe . && GOOS=linux garble build -ldflags '-s -w' -tags phocus -o ./phocus . && GOOS=linux garble build -ldflags '-s -w' -o ./aeacus .

all-dev:
	GOOS=windows go build -tags phocus -o ./phocus.exe . && GOOS=windows go build -o ./aeacus.exe . && GOOS=linux go build -tags phocus -o ./phocus . && GOOS=linux go build -o ./aeacus .

lin:
	GOOS=linux garble build -ldflags '-s -w' -tags phocus -o ./phocus . && GOOS=linux garble build -ldflags '-s -w' -o ./aeacus .

lin-dev:
	GOOS=linux go build -tags phocus -o ./phocus . && GOOS=linux go build -o ./aeacus .

win:
	GOOS=windows garble build -ldflags '-s -w' -tags phocus -o ./phocus.exe . && GOOS=windows garble build -ldflags '-s -w' -o ./aeacus.exe .

win-dev:
	GOOS=windows go build -tags phocus -o ./phocus.exe . && GOOS=windows go build -o ./aeacus.exe .
