.PHONY: build acceptance

build:
	go build -ldflags "-s -w -X 'github.com/friendsofgo/killgrave/internal/app/cmd._version=`git rev-parse --abbrev-ref HEAD`-`git rev-parse --short HEAD`'" -o bin/killgrave cmd/killgrave/main.go

acceptance: build
	@(cd acceptance && go test -count=1 -tags=acceptance -v ./...)