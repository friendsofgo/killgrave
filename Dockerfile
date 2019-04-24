FROM golang:alpine AS build

LABEL MAINTAINER = 'Friends of Go (it@friendsofgo.tech)'

RUN apk add --update git
WORKDIR /go/src/github.com/friendsofgo/killgrave
COPY . .
RUN export GO111MODULE=on && go mod tidy && TAG=$(git describe --tags --abbrev=0) \
    && LDFLAGS=$(echo "-s -w -X main.version="$TAG) \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/killgrave -ldflags "$LDFLAGS" cmd/killgrave/main.go

# Building image with the binary
FROM scratch
COPY --from=build /go/bin/killgrave /go/bin/killgrave
ENTRYPOINT ["/go/bin/killgrave"]