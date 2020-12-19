all: test
	go build -o "foodmap" "./cmd/foodmap/main.go"

run: all
	./foodmap

test:
	go test ./... --count=1
