run:
	go run cmd/*.go -msg "This is test speaking." -lang en

run2:
	go run cmd/*.go -msg "This is test speaking." -lang en -log 2

build:
	go build -i -race -v -o ${GOPATH}/bin/gh ./cmd/

