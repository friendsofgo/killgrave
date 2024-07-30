.PHONY: build
build:
	go build -ldflags "-s -w -X 'github.com/friendsofgo/killgrave/internal/app/cmd._version=`git rev-parse --abbrev-ref HEAD`-`git rev-parse --short HEAD`'" -o bin/killgrave cmd/killgrave/main.go

.PHONY: build-docker
build-docker:
	docker build --build-arg TAG=$(TAG) -t killgrave:$(TAG) .

.PHONY: test
test:
	go test -v -vet=off -race ./...
