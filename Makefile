all: fmt test
	go build "./cmd/foodmap"

fmt:
	go fmt ./...

run: all
	./foodmap

test:
	go test ./... --count=1