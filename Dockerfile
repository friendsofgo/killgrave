FROM golang:1.21-alpine AS build

LABEL MAINTAINER = 'Friends of Go (it@friendsofgo.tech)'

ARG TARGETOS=linux
ARG TARGETARCH=amd64

RUN apk add --update git
RUN apk add ca-certificates
WORKDIR /go/src/github.com/friendsofgo/killgrave
COPY . .
RUN go mod tidy && TAG=$(git describe --tags --abbrev=0) \
    && LDFLAGS=$(echo "-s -w -X github.com/friendsofgo/killgrave/internal/app/cmd._version="docker-$TAG) \
    && CGO_ENABLED=0 GOOS="${TARGETOS}" GOARCH="${TARGETARCH}" go build -a -installsuffix cgo -o /go/bin/killgrave -ldflags "$LDFLAGS" cmd/killgrave/main.go

# Building image with the binary
FROM scratch
COPY --from=build /go/bin/killgrave /go/bin/killgrave
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/go/bin/killgrave"]