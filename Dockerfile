FROM node:16-alpine AS node

FROM golang:1.19-alpine AS build

LABEL MAINTAINER = 'Friends of Go (it@friendsofgo.tech)'

# Copy Node binaries
COPY --from=node /usr/lib /usr/lib
COPY --from=node /usr/local/share /usr/local/share
COPY --from=node /usr/local/lib /usr/local/lib
COPY --from=node /usr/local/include /usr/local/include
COPY --from=node /usr/local/bin /usr/local/bin

# Install other dependencies
RUN apk add --update git yarn
RUN apk add ca-certificates

# Copy source code
WORKDIR /go/src/github.com/friendsofgo/killgrave
COPY . .

# Build Node app
WORKDIR /go/src/github.com/friendsofgo/killgrave/debugger
RUN yarn install --immutable && yarn build

# Build Go binary
WORKDIR /go/src/github.com/friendsofgo/killgrave
RUN go mod tidy && TAG=$(git describe --tags --abbrev=0) \
    && LDFLAGS=$(echo "-s -w -X main.version="$TAG) \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/killgrave -ldflags "$LDFLAGS" cmd/killgrave/main.go

# Building image with the binary
FROM scratch
COPY --from=build /go/bin/killgrave /go/bin/killgrave
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/go/bin/killgrave"]