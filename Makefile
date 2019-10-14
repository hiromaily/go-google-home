
###############################################################################
# Managing Dependencies
###############################################################################
.PHONY: update
update:
	GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	go get -u -d -v ./...


###############################################################################
# Golang formatter and detection
###############################################################################
.PHONY: lint
lint:
	golangci-lint run --fix

.PHONY: imports
imports:
	./scripts/imports.sh

###############################################################################
# Build
###############################################################################
.PHONY: build
build:
	go build -v -o ${GOPATH}/bin/gh ./cmd/

.PHONY: build-mac
build-mac: GOOS=darwin GOARCH=amd64
build-mac:
	go build -v -o ./cmd/releases/darwin_amd64/gh ./cmd/

.PHONY: build-linux
build-linux: GOOS=linux
build-linux:
	#GOOS=linux go build -v -o ./releases/linux_amd64/gh ./cmd/
	go build -v -o ./cmd/releases/linux_amd64/gh ./cmd/

.PHONY: build-linux-arm
build-linux-arm: GOOS=linux GOARCH=arm GOARM=5
build-linux-arm:
	#GOOS=linux GOARCH=arm GOARM=5 go build -v -o ./releases/linux_arm/gh ./cmd/
	go build -v -o ./cmd/releases/linux_arm/gh ./cmd/

.PHONY: release-all
release-all: build-linux build-linux-arm build-mac


###############################################################################
# Execute
###############################################################################
# Sample for saying something in English.
.PHONY: say-en
say-en:
	gh -msg "Thank you."

# Sample for saying something in Japanese.
.PHONY: say-ja
say-ja:
	gh  -msg "ありがとうございます" -lang ja

# Sample for saying something in Dutch.
.PHONY: say-nl
say-nl:
	gh  -msg "Dank je" -lang nl

# Sample for saying something in German.
.PHONY: say-de
say-de:
	gh  -msg "Danke." -lang de

# Sample for saying something in French.
.PHONY: say-fr
say-fr:
	gh -msg "Merci." -lang fr

# Sample for saying by specific sound volume.
.PHONY: say-en2
say-en2:
	gh -msg "Thank you." -vol 0.3


# Sample for playing music.
.PHONY: play-music
play-music:
	gh -music "https://raw.githubusercontent.com/hiromaily/go-google-home/master/asetts/music/bensound-dubstep.mp3"


# Sample for using specific IP address of Google Home.
.PHONY: play-music
say-with_address:
	gh -msg "It reaches to specific IP address." -addr "10.0.0.1:8009"


# Sample for saying something in English with `debug` log.
.PHONY: say-debug
say-debug:
	gh -msg "This displays debug log." -log 1


###############################################################################
# server mode
###############################################################################
# Sample for server mode.
.PHONY: server
server:
	gh -server

# Sample to post message to server by HTTPie
.PHONY: post-msg
post-msg:
	http POST http://localhost:8080/speak text="It's sunny day today."
