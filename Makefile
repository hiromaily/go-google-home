# Build
build:
	go build -i -race -v -o ${GOPATH}/bin/gh ./cmd/


# Sample for saying something in English.
say-en:
	#go run cmd/*.go -msg "Thank you." -lang en
	gh -msg "Thank you." -lang en

# Sample for saying something in Japanese.
say-ja:
	#go run cmd/*.go -msg "ありがとうございます" -lang ja
	gh  -msg "ありがとうございます" -lang ja

# Sample for saying something in French.
say-fr:
	#go run cmd/*.go -msg "Merci." -lang fr
	gh -msg "Merci." -lang fr

# Sample for saying something in German.
say-de:
	#go run cmd/*.go -msg "Danke." -lang de
	gh  -msg "Danke." -lang de


# Sample for saying something in English with `debug` log.
say-debug:
	#go run cmd/*.go -msg "This displays debug log." -log 1
	gh  -msg "This displays debug log." -log 1

# Sample for server mode.
server:
	#go run cmd/*.go -server
	gh -server

# Sample to post message to server by HTTPie
post-msg:
	http POST http://localhost:8080/speak text="It's sunny day today."
