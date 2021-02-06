
CURRENTDIR=`pwd`
modVer=$(shell cat go.mod | head -n 3 | tail -n 1 | awk '{print $2}' | cut -d'.' -f2)
currentVer=$(shell go version | awk '{print $3}' | sed -e "s/go//" | cut -d'.' -f2)
gitTag=$(shell git tag | head -n 1)

###############################################################################
# Managing Dependencies
###############################################################################
.PHONY: check-ver
check-ver:
	#echo $(modVer)
	#echo $(currentVer)
	@if [ ${currentVer} -lt ${modVer} ]; then\
		echo go version ${modVer}++ is required but your go version is ${currentVer};\
	fi

.PHONY: update
update:
	GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	go get -u -d -v ./...


###############################################################################
# Golang formatter and detection
###############################################################################
.PHONY: imports
imports:
	./scripts/imports.sh

.PHONY: lint
lint:
	golangci-lint run

.PHONY: lintfix
lintfix:
	golangci-lint run --fix

.PHONY: lintall
lintall: imports lint


###############################################################################
# Build
###############################################################################
.PHONY: build
build:
	go build -v -o ${GOPATH}/bin/gh ./cmd/gh/

.PHONY: build-version
build-version:
	go build -ldflags "-X main.version=${gitTag}" -v -o ${GOPATH}/bin/gh ./cmd/gh/

.PHONY: run
run:
	go run -v ./cmd/gh/ speak -msg "My name is Hiroki. Nice to meet you. How's it going? Thank you."


###############################################################################
# Execute
###############################################################################
# Sample for saying something in English.
.PHONY: say-en
say-en:
	gh speak -msg "Thank you."

# Sample for saying something in Japanese.
.PHONY: say-ja
say-ja:
	gh speak -msg "ありがとうございます" -lang ja

# Sample for saying something in Dutch.
.PHONY: say-nl
say-nl:
	gh -lang nl speak -msg "Dank je"

# Sample for saying something in German.
.PHONY: say-de
say-de:
	gh -lang de speak -msg "Danke."

# Sample for saying something in French.
.PHONY: say-fr
say-fr:
	gh -lang fr -msg "Merci."

# Sample for saying by specific sound volume.
.PHONY: say-en2
say-en2:
	gh -vol 0.3 speak -msg "Thank you."


# Sample for playing music.
.PHONY: play-music
play-music:
	gh music -file "https://raw.githubusercontent.com/hiromaily/go-google-home/master/asetts/music/bensound-dubstep.mp3"


# Sample for using specific IP address of Google Home.
.PHONY: say-with-address
say-with-address:
	gh -addr "10.0.0.1:8009" -msg "It reaches to specific IP address."


###############################################################################
# server mode
###############################################################################
# Sample for server mode.
.PHONY: server
server:
	gh server

# Sample to post message to server by HTTPie
.PHONY: post-msg
post-msg:
	http POST http://localhost:8080/speak text="It's sunny day today."
