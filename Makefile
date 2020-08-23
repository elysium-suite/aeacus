.DEFAULT_GOAL := all

all:
	GOOS=windows garble build -ldflags '-s -w' -tags phocus -o ./phocus.exe . && \
	GOOS=windows garble build -ldflags '-s -w' -o ./aeacus.exe . && \
	echo "Windows production build successful!" && \
	GOOS=linux garble build -ldflags '-s -w' -tags phocus -o ./phocus . && \
	GOOS=linux garble build -ldflags '-s -w' -o ./aeacus . && \
	echo "Linux production build successful!"

all-dev:
	GOOS=windows go build -tags phocus -o ./phocus.exe . && \
	GOOS=windows go build -o ./aeacus.exe . && \
	echo "Windows development build successful!" && \
	GOOS=linux go build -tags phocus -o ./phocus . && \
	GOOS=linux go build -o ./aeacus . && \
	echo "Linux development build successful!"

lin:
	GOOS=linux garble build -ldflags '-s -w' -tags phocus -o ./phocus . && \
	GOOS=linux garble build -ldflags '-s -w' -o ./aeacus . && \
	echo "Linux production build successful!"

lin-dev:
	GOOS=linux go build -tags phocus -o ./phocus . && \
	GOOS=linux go build -o ./aeacus . && \
	echo "Linux development build successful!"

win:
	GOOS=windows garble build -ldflags '-s -w' -tags phocus -o ./phocus.exe . && \
	GOOS=windows garble build -ldflags '-s -w' -o ./aeacus.exe . && \
	echo "Windows production build successful!"

win-dev:
	GOOS=windows go build -tags phocus -o ./phocus.exe . && GOOS=windows go build -o ./aeacus.exe . && \
	echo "Windows development build successful!"
