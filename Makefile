# Build
build:
	go build -i -race -v -o ${GOPATH}/bin/gh ./cmd/

# Sample for saying something in English.
run:
	go run cmd/*.go -msg "This is test speaking." -lang en

# Sample for saying something in English without debug log.
run2:
	go run cmd/*.go -msg "This is test speaking." -lang en -log 2

# Sample for server mode.
server:
	go run cmd/*.go -server

# Sample to post message to server by HTTPie
say:
	http POST http://localhost:8080/speak text="It's sunny day today."


