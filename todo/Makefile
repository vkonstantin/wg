build:
	go build

test:
	go test ./... -v -timeout 90s

docker:
	docker rmi -f wg/todo-image
	docker build -t wg/todo-image .

.PHONY: build test docker

