test:
	go test ./... --count=1

build:
	go build -o "foodmap" "./cmd/foodmap/main.go"
