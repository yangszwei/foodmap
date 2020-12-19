all: test
	go build -o "foodmap" "./cmd/foodmap/main.go"

fmt:
	go fmt ./...

run: all
	./foodmap

test:
	go test ./... --count=1
