# Build
build:
	go build -i -race -v -o ${GOPATH}/bin/gh ./cmd/

build-mac: GOOS=darwin GOARCH=amd64
build-mac:
	go build -i -race -v -o ./cmd/releases/darwin_amd64/gh ./cmd/

build-linux: GOOS=linux
build-linux:
	#GOOS=linux go build -v -o ./releases/linux_amd64/gh ./cmd/
	go build -v -o ./cmd/releases/linux_amd64/gh ./cmd/

build-linux-arm: GOOS=linux GOARCH=arm GOARM=5
build-linux-arm:
	#GOOS=linux GOARCH=arm GOARM=5 go build -v -o ./releases/linux_arm/gh ./cmd/
	go build -v -o ./cmd/releases/linux_arm/gh ./cmd/

release-all: build-linux build-linux-arm build-mac


# Sample for saying something in English.
say-en:
	gh -msg "Thank you."

# Sample for saying something in Japanese.
say-ja:
	gh  -msg "ありがとうございます" -lang ja

# Sample for saying something in Dutch.
say-nl:
	gh  -msg "Dank je" -lang nl

# Sample for saying something in German.
say-de:
	gh  -msg "Danke." -lang de

# Sample for saying something in French.
say-fr:
	gh -msg "Merci." -lang fr


# Sample for using specific IP address of Google Home.
say-with_address:
	gh  -msg "It reaches to specific IP address." -addr "10.0.0.1:8009"


# Sample for saying something in English with `debug` log.
say-debug:
	gh  -msg "This displays debug log." -log 1


# Sample for server mode.
server:
	gh -server

# Sample to post message to server by HTTPie
post-msg:
	http POST http://localhost:8080/speak text="It's sunny day today."
