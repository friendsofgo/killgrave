[![CircleCI](https://circleci.com/gh/friendsofgo/killgrave/tree/master.svg?style=svg)](https://circleci.com/gh/friendsofgo/killgrave/tree/master)
[![Version](https://img.shields.io/github/release/friendsofgo/killgrave.svg?style=flat-square)](https://github.com/friendsofgo/killgrave/releases/latest)
[![codecov](https://codecov.io/gh/friendsofgo/killgrave/branch/master/graph/badge.svg)](https://codecov.io/gh/friendsofgo/killgrave)
[![Go Report Card](https://goreportcard.com/badge/github.com/friendsofgo/killgrave)](https://goreportcard.com/report/github.com/friendsofgo/killgrave)
[![GoDoc](https://godoc.org/graphql.co/graphql?status.svg)](https://godoc.org/github.com/friendsofgo/killgrave)
[![FriendsOfGo](https://img.shields.io/badge/powered%20by-Friends%20of%20Go-73D7E2.svg)](https://friendsofgo.tech)

<p align="center">
  <img src="https://res.cloudinary.com/fogo/image/upload/c_scale,w_350/v1555701634/fogo/projects/gopher-killgrave.png" alt="Golang Killgrave"/>
</p>

# Killgrave

Killgrave is a simulator for HTTP-based APIs, in simple words a **Mock Server**, very easy to use made in **Go**.

## Getting started
Install `killgrave` using go:

```sh
$ GO111MODULE=off go get -u github.com/friendsofgo/killgrave/cmd/killgrave
```

Install `killgrave` using [homebrew](https://brew.sh/index_es):

```sh
$ brew install friendsofgo/tap/killgrave
```

Or you can download the binary for your arch on:

[https://github.com/friendsofgo/killgrave/releases](https://github.com/friendsofgo/killgrave/releases)

## Using Killgrave

Use `killgrave` with default flags:

```sh
$ killgrave
2019/04/14 23:53:26 The fake server is on tap now: http://localhost:3000
```
Or custome your server with this flags:
```sh
 -host string
        if you run your server on a different host (default "localhost")
 -imposters string
        directory where your imposter are saved (default "imposters")
 -port int
        por to run the server (default 3000)
 -version
        show the version of the application
```

## How to use

### Create an imposter
You must be create an imposter to start to use the application

```json
imposters/create_gopher.json

{
    "request": {
        "method": "POST",
        "endpoint": "/gophers",
        "schema_file": "schemas/create_gopher_request.json",
        "headers": {
            "Content-Type": [
                "application/json"
            ]
        }
    },
    "response": {
        "status": 200,
        "headers": {
            "Content-Type": [
                "application/json"
            ]
        },
        "bodyFile": "responses/create_gopher_response.json"
    }
}
```
And its related files

```json
schemas/create_gopher_request.json
{
    "type": "object",
    "properties": {
        "data": {
            "type": "object",
            "properties": {
                "type": {
                    "type": "string",
                    "enum": [
                        "gophers"
                    ]
                },
                "attributes": {
                    "type": "object",
                    "properties": {
                        "name": {
                            "type": "string"
                        },
                        "color": {
                            "type": "string"
                        },
                        "age": {
                            "type": "integer"
                        }
                    },
                    "required": [
                        "name",
                        "color",
                        "age"
                    ]
                }
            },
            "required": [
                "type",
                "attributes"
            ]
        }
    },
    "required": [
        "data"
    ]
}
```

```json
responses/create_gopher_response.json
{
    "data": {
        "type": "gophers",
        "id": "01D8EMQ185CA8PRGE20DKZTGSR",
        "attributes": {
            "name": "Zebediah",
            "color": "Purple",
            "age": 55
        }
    }
}
```

And then with the server on tap you can execute your request:
```sh
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{
            "data": {
                "type": "gophers",
                "attributes": {
                "name": "Zebediah",
                "color": "Purple",
                "age": 55
                }
            }
    }' \
  http://localhost:3000/gophers
```

### Docker

The application is also available through [Docker](https://hub.docker.com/r/friendsofgo/killgrave), just run:

```bash
docker run -it --rm -p 3000:3000 friendsofgo/killgrave
```

Remember to use the [-p](https://docs.docker.com/engine/reference/run/) flag to expose the container port where the application is listening (3000 by default).

NOTE: If you want to use `killgrave` through Docker at the same time you use your own dockerised HTTP-based API, be careful with networking issues.

## Features
* Imposters created in json
* Validate json schemas on requests
* Validate requests headers
* Check response status
* All content-type bodies
* Write body files (XML, JSON, HTML...)
* Write bodies in line
* Regex for using on endpoint urls
* Allow write headers on response
* Allow imposter's matching by request schema
* Dynamic responses based on regex endpoint or request schema

## Next Features
- [ ] Dynamic responses based on headers
- [ ] Dynamic responses based on query params
- [ ] Allow write multiples imposters by file
- [ ] Proxy server
- [ ] Record proxy server
- [ ] Better documentation with examples of each feature

## Contributing
[Contributions](https://github.com/friendsofgo/killgrave/issues?q=is%3Aissue+is%3Aopen) are more than welcome, if you are interested please fork this repo and send your Pull Request.

## License
MIT License, see [LICENSE](https://github.com/friendsofgo/killgrave/blob/master/LICENSE)
