[![CircleCI](https://circleci.com/gh/friendsofgo/killgrave/tree/master.svg?style=svg)](https://circleci.com/gh/friendsofgo/killgrave/tree/master)
[![Version](https://img.shields.io/github/release/friendsofgo/killgrave.svg?style=flat-square)](https://github.com/friendsofgo/killgrave/releases/latest)
[![codecov](https://codecov.io/gh/friendsofgo/killgrave/branch/master/graph/badge.svg)](https://codecov.io/gh/friendsofgo/killgrave)
[![Go Report Card](https://goreportcard.com/badge/github.com/friendsofgo/killgrave)](https://goreportcard.com/report/github.com/friendsofgo/killgrave)
[![FriendsOfGo](https://img.shields.io/badge/powered%20by-Friends%20of%20Go-73D7E2.svg)](https://friendsofgo.tech)

</br>
<a href="https://www.buymeacoffee.com/friendsofgo" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: auto !important;width: 100px !important;" ></a>

<p align="center">
  <img src="https://res.cloudinary.com/fogo/image/upload/c_scale,w_350/v1555701634/fogo/projects/gopher-killgrave.png" alt="Golang Killgrave"/>
</p>

# Killgrave

Killgrave is a simulator for HTTP-based APIs, in simple words a **Mock Server**, very easy to use made in **Go**.

## Versions

Actually `Killgrave` working using `sem ver` but in 0 version, for that the 'minor' will be change when some broken changes are added
to the application, and the 'patch' would be changed when a new feature with new changes are added or for bug fixing, when 
v1.0.0 will be released, `Killgrave` will be change to use a normal use of `sem ver`. 

## Branch master

Master branch contains all the latest changes on the application, for that the branch `master` is not an stable version.
If you want to have those changes then you will need to use the branch master, but remember that some breaking changes can be added
to this branch.

## Future releases
This project is on continuos improvement so you can check the [Issues](https://github.com/friendsofgo/killgrave/issues) created for give it support or create your owns.
Furthermore on the `branch master` you can find the [CHANGELOG.md](https://github.com/friendsofgo/killgrave/blob/master/CHANGELOG.md) file that contains all the new feature would be added
on the next release. 

## Getting started

Install `killgrave` using go:

```sh
$ go get -u github.com/friendsofgo/killgrave/cmd/killgrave@{version}
```

`version` must be substituted by the `version` that you want to install, otherwise master would be installed.

Install `killgrave` using [homebrew](https://brew.sh/index_es):

```sh
$ brew install friendsofgo/tap/killgrave
```

Or you can download the binary for your arch on:

[https://github.com/friendsofgo/killgrave/releases](https://github.com/friendsofgo/killgrave/releases)

### Docker

The application is also available through [Docker](https://hub.docker.com/r/friendsofgo/killgrave), just run:

```bash
docker run -it --rm -p 3000:3000 -v $PWD/:/home -w /home friendsofgo/killgrave
```
Remember to use the [-p](https://docs.docker.com/engine/reference/run/) flag to expose the container port where the application is listening (3000 by default).

NOTE: If you want to use `killgrave` through Docker at the same time you use your own dockerised HTTP-based API, be careful with networking issues.

## Using Killgrave

Use `killgrave` with default flags:

```sh
$ killgrave
2019/04/14 23:53:26 The fake server is on tap now: http://localhost:3000
```
Or custome your server with this flags:
```sh
  -config string
        path with configuration file
  -host string
        if you run your server on a different host (default "localhost")
  -imposters string
        directory where your imposters are saved (default "imposters")
  -port int
        port to run the server (default 3000)
  -proxy-mode string
        proxy mode you can choose between (all, missing or none) (default "none")
  -proxy-url string
        proxy url, you need to choose a proxy-mode
  -version
        show the _version of the application
  -watcher
        file watcher, reload the server with each file change
```

Use `killgrave` with config file:

First of all you need create a file with a valid config, i.e:

```yaml
#config.yml

imposters_path: "imposters"
port: 3000
host: "localhost"
proxy:
  url: https://example.com
  mode: missing
cors:
  methods: ["GET"]
  headers: ["Content-Type"]
  exposed_headers: ["Cache-Control"]
  origins: ["*"]
  allow_credentials: true
```

Historically, the options `imposters_path`, `port`, `host` were mandatory when using a configuration file.

However, since last versions, they are no longer needed, so you can simply override those options you want to.
Furthermore the `imposters_path` option in previous version towards reference to the path where the app was launched, but
in the last version the `imposters_path` option is relative on where the config file is.

The `delay` option is a time that server waits before response. This can help simulate network issues, or server high load. You must write delay as string with postfix indicating time unit (see [this](https://golang.org/pkg/time/#ParseDuration) for more info about actual format). Also you can specify minimum and maximum delays using separator ':', server respond delay will be choosen at random between this values.

Default value is no delay at all.

The option `cors` still being optional and his options can be an empty array,

If you want more information about the CORS options, visit the [CORS section](#CORS).

## How to use

### Create an imposter
You must be create an imposter to start to use the application, only files with the `.imp.json` extension will be interpreted as imposters files, and the base path for
the rest of the files will be the path of the `.imp.json` file.

You need to organize your imposters from more restrictive to less. We use a rule-based system for create each imposter, for this reason you need to organize your imposters in the way that more restrictive to less, like the example below.

```json
imposters/create_gopher.imp.json

[
    {
        "request": {
            "method": "POST",
            "endpoint": "/gophers",
            "schemaFile": "schemas/create_gopher_request.json",
            "headers": {
                "Content-Type": "application/json",
                "Return-Error": "error"
            }
        },
        "response": {
            "status": 500,
            "headers": {
                "Content-Type": "application/json"
            },
            "body": "{\"error\": \"Server error ocurred\"}",
            "delay": "1s:5s"
        }
    },
    {
        "request": {
            "method": "POST",
            "endpoint": "/gophers",
            "schemaFile": "schemas/create_gopher_request.json",
            "headers": {
                "Content-Type": "application/json"
            }
        },
        "response": {
            "status": 200,
            "headers": {
                "Content-Type": "application/json"
            },
            "bodyFile": "responses/create_gopher_response.json"
        }
    }
]
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

## CORS

If you want to use `killgrave` on your client application you must consider to configure correctly all about CORS, thus we offer the possibility to configure as you need through a config file.

In the CORS section of the file you can find the next options:

- **methods** (string array)
  
  Represent the **Access-Control-Request-Method header**, if you don't specify it or if you do leave it as any empty array, the default values will be:

  `"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE", "PATCH", "TRACE", "CONNECT"`

- **headers** (string array)
  
  Represent the **Access-Control-Request-Headers header**, if you don't specify it or if you do leave it as any empty array, the default values will be:

  `"X-Requested-With", "Content-Type", "Authorization"`

- **exposed_headers** (string array)
  
  Represent the **Access-Control-Expose-Headers header**, if you don't specify it or if you do leave it as any empty array, the default values will be:

  `"Cache-Control", "Content-Language", "Content-Type", "Expires", "Last-Modified", "Pragma"`

- **origins** (string array)
  
  Represent the **Access-Control-Allow-Origin header**, if you don't specify or leave as empty array this options has not default value

- **allow_credentials** (boolean)
  
  Represent the **Access-Control-Allow-Credentials header** you must indicate if true or false
  
## Proxy
You can use Killgrave with a proxy mode, what does that mean, well you can use the flags `proxy-mode` and `proxy-url` or the configuration file to declare one of these three modes:
* none: by default mode, with this mode you don't use any proxy, and the mock server always will looks into the files with `.imp.json` extensions.
* missing: with this mode, the mock server will looks into the files with `.imp.json` extensions but if you call to an endpoint that doesn't exists then the mock server will call to the real server, declaring on the `proxy-url`.
* all: the mock server only will call to the real server, declaring on the `proxy-url`.

The `proxy-url` must be the root path. For example if we have endpoint api like, `http://example.com/things`, the `proxy-url` will be, `http://example.com`

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
* Dynamic responses based on headers
* Dynamic responses based on query params
* Allow organize your imposters with structured folders
* Allow write multiple imposters by file
* Run mock server with predefined configuration with config yaml file
* Configure your CORS server options
* Simulate network issues and server high loads with imposter repsonse delay
* Proxy server

## Next Features
- [ ] Record proxy server
- [ ] Better documentation with examples of each feature
- [ ] Validate request body XML

## Contributing
[Contributions](https://github.com/friendsofgo/killgrave/issues?q=is%3Aissue+is%3Aopen) are more than welcome, if you are interested please fork this repo and send your Pull Request.

## License
MIT License, see [LICENSE](https://github.com/friendsofgo/killgrave/blob/master/LICENSE)
