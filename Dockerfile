FROM --platform=$BUILDPLATFORM golang:1.27-alpine AS build

LABEL maintainer="Friends of Go (it@friendsofgo.tech)"

RUN apk add --no-cache ca-certificates
WORKDIR /go/src/github.com/friendsofgo/killgrave
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS="${TARGETOS}" GOARCH="${TARGETARCH}" go build \
    -buildvcs=false -trimpath \
    -ldflags "-s -w -X github.com/friendsofgo/killgrave/internal/app/cmd._version=docker-${VERSION}" \
    -o /go/bin/killgrave ./cmd/killgrave

FROM scratch
COPY --from=build /go/bin/killgrave /go/bin/killgrave
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/go/bin/killgrave"]
