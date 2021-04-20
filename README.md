<p align="center">
  <img src="https://res.cloudinary.com/fogo/image/upload/c_scale,w_350/v1555701634/fogo/projects/gopher-killgrave.png" alt="Golang Killgrave"/>
</p>

# Killgrave

Killgrave is a simulator for HTTP-based APIs, in simple words a **Mock Server**, very easy to use, made in **Go**.

[![CircleCI](https://circleci.com/gh/friendsofgo/killgrave/tree/master.svg?style=svg)](https://circleci.com/gh/friendsofgo/killgrave/tree/master)
[![Version](https://img.shields.io/github/release/friendsofgo/killgrave.svg?style=flat-square)](https://github.com/friendsofgo/killgrave/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/friendsofgo/killgrave)](https://goreportcard.com/report/github.com/friendsofgo/killgrave)
[![Total alerts](https://img.shields.io/lgtm/alerts/g/friendsofgo/killgrave.svg?logo=lgtm&logoWidth=18)](https://lgtm.com/projects/g/friendsofgo/killgrave/alerts/)
[![FriendsOfGo](https://img.shields.io/badge/powered%20by-Friends%20of%20Go-73D7E2.svg)](https://friendsofgo.tech)

<p>
<a href="https://www.buymeacoffee.com/friendsofgo" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: auto !important;width: 100px !important;" ></a>
</p>

# Table of Content
- [Overview](#overview)
- [Concepts](#concepts)
    * [Imposters](#imposters)
- [Installing](#installing)
    * [Go Toolchain](#go-toolchain)
    * [Homebrew](#homebrew)
    * [Docker](#docker)
    * [Other](#other)
- [Getting Started](#getting-started)
    * [Using Killgrave by command line](#using-killgrave-by-command-line)
    * [Using Killgrave by config file](#using-killgrave-by-config-file)
    * [Configure CORS](#configure-cors)
    * [Prepare Killgrave for Proxy Mode](#prepare-killgrave-for-proxy-mode)
    * [Create an Imposter](#create-an-imposter)
    * [Imposters structure](#imposters-structure)
    * [Create an Imposter using regex](#create-an-imposter-using-regex)
    * [Create an imposter using JSON Schema](#create-an-imposter-using-json-schema)
    * [Create an imposter with delay](#create-an-imposter-with-delay)
    * [Create an imposter with dynamic responses](#create-an-imposter-with-dynamic-responses)
- [Contributing](#contributing)
- [License](#license)

## Overview

**Killgrave** is a tool providing a simple way to create a powerful simulator for HTTP-based APIs.

The Killgrave's philosophy is provides an easy way to configure your *mock server*, trying always that you don't need learn 
another tooling language, for that reason we use `json` and `yaml` to generate all necessary configurations.

Killgrave provides:

* Easy way to create imposters files, using `json`
* The possibility to validate your owns [json schemas](https://json-schema.org/) on requests.
* Validations the requests headers.
* Using regex to allow different parameters on headers and urls.
* Custom body and dynamic body responses
* Using all content-types bodies, (`application/json`, `text/html`, `text/plain`, `image/jpeg`, etc. )
* Configure the headers on body responses.
* Configure your CORS server options.
* Simulate network issues and server high loads with imposter responses delay.
* Using configurable proxy server to call to the original server.
* Run the tool using flag or using a config file.
* Run your mock server using a watcher to not reload on each change.

## Concepts

### Imposters
Maybe the most common concept inside of Killgrave tool is the `imposters`. Imposters are the files that we will use to configure
how we want that our server respond.

You can identify an `imposter` file on Killgrave because they must have the extension `.imp.json`.

Yo can learn more about how to configure the imposters on the [Imposter Configuration Section](#imposter).

## Installing
> :warning:  Killgrave is a very robust tool and is using on some companies even on production environments, but we don't have a
version 1 for today. For that Killgrave works using [`SemVer`](https://semver.org/) but in 0 version, which means that the 'minor' will be changed when some broken changes are introduced into the application, and the 'patch' will be changed when a new feature with new changes is added or for bug fixing. As soon as v1.0.0 be released, `Killgrave` will start to use [`SemVer`](https://semver.org/) as usual.

You can install Killgrave on different ways, but all of them are very simple:

### Go Toolchain

One of them is of course using `go get`, Killgrave is a Go project so could be compiled using the `go toolchain`:

```sh
$ go get -u github.com/friendsofgo/killgrave/cmd/killgrave@{version}
```

`version` must be substituted by the `version` that you want to install, otherwise master would be installed.

### Homebrew 

If you are a mac user, you could install Killgrave using, [Homebrew](https://brew.sh/):

```sh
$ brew install friendsofgo/tap/killgrave
```

:warning:  If you are installing via homebrew, you only get the [last Killgrave version](https://github.com/friendsofgo/killgrave/releases), we hope fix this soon.

### Docker

The application is also available through [Docker](https://hub.docker.com/r/friendsofgo/killgrave).

```bash
docker run -it --rm -p 3000:3000 -v $PWD/:/home -w /home friendsofgo/killgrave -h 0.0.0.0
```

`-p 3000:3000` [publishes](https://docs.docker.com/engine/reference/run/#expose-incoming-ports) port 3000 inside the
container (where Killgrave is listening by default) to port 3000 on the host machine.

`-h 0.0.0.0` is necessary to allow Killgrave to listen and respond to requests from outside the container (the default,
`localhost`, will not capture requests from the host network).

### Other

Windows and Linux users can download binaries from the [Github Releases](https://github.com/friendsofgo/killgrave/releases) page.

## Getting Started

The very easy way to start with Killgrave is just execute at it is:

```sh
$ killgrave
```

While you are welcome to provide your own configuration, Killgrave will follow the next default configuration:

* **imposters path**: `imposters`
* **host**: `localhost`
* **port**: `3000`
* **CORS**: `[]`
* **proxy**: `none`
* **watcher**: `false`

### Using Killgrave by command line

As we see below Killgrave have a default configuration to start using out of the box, but you could change all of this using
his command tool, except for the `CORS` that you only will can configure using the [config file](#using-killgrave-by-config-file).

```sh
$ killgrave -h

  -config string
        path with configuration file
  -host string
        if you run your server on a different host (default "localhost")
  -imposters string
        directory where your imposters are saved (default "imposters")
  -port int
        port to run the server (default 3000)
  -secure bool
        if you run your server using TLS (https)
  -proxy-mode string
        proxy mode you can choose between (all, missing or none) (default "none")
  -proxy-url string
        proxy url, you need to choose a proxy-mode
  -version
        show the _version of the application
  -watcher
        file watcher, reload the server with each file change
```

### Using Killgrave by config file

If we want a more permanent configuration, we could use the option `-config` to define where is our config file.

The config file must be a `yml file`, his structure should be like this:

```yaml
#config.yml

imposters_path: "imposters"
port: 3000
host: "localhost"
proxy:
  url: https://example.com
  mode: missing
watcher: true
cors:
  methods: ["GET"]
  headers: ["Content-Type"]
  exposed_headers: ["Cache-Control"]
  origins: ["*"]
  allow_credentials: true
watcher: true
secure: true
```

As you can see, you could configure all the options on a very easy way. Keep in mind that the routes are based on the config file are.

example:

```
mymock/
    imposters/
        config.yml
        swapi_people.imp.json
        swapi_planets.imp.json
    Dockerfile
    MakeFile
```

Then on your config file, you will need configure `imposters` option like `.` because the `imposters file` are located in the
same folder.

Historically, the options `imposters_path`, `port`, `host` were mandatory when using a configuration file.

However, since the last version, they are no longer needed, so you can simply override those options if you want to.
Furthermore, the `imposters_path` option in previous version towards reference to the path where the app was launched, but in the last version it is relative on where the config file is.

The option `cors` still being optional and its options can be an empty array.
If you want more information about the CORS options, visit the [CORS section](#configure-cors).

The `watcher` configuration field is optional, with this setting you can enable hot-reloads on imposter changes. Disabled by default.

The `secure` configuration field is optional, with this setting you can run your server using TLS options with a dummy certificate, so to make it work with the `HTTPS` protocol. Disabled by default.

The option `proxy` allow to configure the mock on a proxy mode, that means that you could configure an fallback urls, for all
your calls or only the missing ones or none. More information: [Proxy Section](#prepare-killgrave-for-proxy-mode)

## How to use

### Configure CORS

If you want to use `killgrave` on your client application you must consider to configure correctly all about CORS, thus we offer the possibility to configure it as you need through a config file.

In the CORS section of the file you can find the following options:

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
  
  Represent the **Access-Control-Allow-Origin header**, if you don't specify or leave as empty array this option has not default value

- **allow_credentials** (boolean)
  
  Represent the **Access-Control-Allow-Credentials header** you must indicate if true or false

### Prepare Killgrave for Proxy Mode

You can use Killgrave with a proxy mode which means that you can use the flags `proxy-mode` and `proxy-url` or the configuration file to declare one of these three modes:
* none: by default mode, with this mode you don't use any proxy, and the mock server will always look into the files with `.imp.json` extension.
* missing: with this mode, the mock server will look into the files with `.imp.json` extension, but if you call to an endpoint that doesn't exist, then the mock server will call to the real server, declared on the `proxy-url` configuration variable.
* all: the mock server only will call to the real server, declared on the `proxy-url` configuration variable.

The `proxy-url` must be the root path. For example, if we have endpoint api like, `http://example.com/things`, the `proxy-url` will be, `http://example.com`

### Create an Imposter

You must be creating an imposter to start using the application, only files with the `.imp.json` extension will be interpreted as imposter files, and the base path for the rest of the files will be the path of the `.imp.json` file.

You need to organize your imposters from more restrictive to less. We use a rule-based system for creating each imposter, for this reason you need to organize your imposters from the most restrictive to the least, like the example below.

imposter example:

```json
[
    {
        "request": {
            "method": "GET",
            "endpoint": "/gophers/01D8EMQ185CA8PRGE20DKZTGSR"            
        },
        "response": {
            "status": 200,
            "headers": {
                "Content-Type": "application/json"
            },
            "body": "{\"data\":{\"type\":\"gophers\",\"id\":\"01D8EMQ185CA8PRGE20DKZTGSR\",\"attributes\":{\"name\":\"Zebediah\",\"color\":\"Purples\",\"age\":55}}}"
        }
    }
]
```

This a very simple example but Killgrave have more possibilities to configure your imposters, let see some of them on the
next sections.

:warning:  Remember that you will need escape any special char, in the properties that admit text.

### Imposters Structure

The imposter object can be divided on two parts:

* [Request](#request)
* [Response](#response)

#### Request

For configure our imposter we will need to declare a `request`, the `request` is the object which allow the Killgrave
engine identify with which endpoint should match, and have the next properties:

* `method` (<span style="color:red">mandatory</span>): with this property we will define the method of our request, 
you could configure whatever verb on the [http protocol](https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods) was allowed. 
* `endpoint` (<span style="color:red">mandatory</span>): it is a relative path with the endpoint that you want to mock. Admit regex to configure
dynamic params.
* `schemaFile`: if you want that the request could be validated using a json schema, you will need to indicate the path where your
json schema are.
* `params: if you want to specify a restrictive query param, you can add a list of them, more information: 
[Create an imposter with query params](#create-an-imposter-with-query-params), this option admit regex.
* `headers`: if you want to specify a restrictive headers, like authorization or any other, you can define a list
of headers that you request will need to match, more info: [Create an imposter with headers](#create-an-imposter-with-headers).

#### Response

On the other hand we have the `response` object, the `response` object will be in charge to generate the output that
we want, related with the `request` that was executed.

The response, allow the next properties:

* `status` (<span style="color:red">mandatory</span>): this will be the status that the response will return.
* `body` or `bodyFile`: if you want that your response return any kind result, you need to define `body` or `bodyFile` param,
the difference between them, is that `body` allow writing directly the output, and `bodyFile`, is a route to the file with
the output, this is very useful in case of a large output. You can also remove, this property if you want not return any content.
* `headers`: sometime we need a specific header on our `response`, for example if you want to specify a special `content-type``
so with this property you can define which headers the response will return.
* `delay`: is the time the server waits before responding. This can help simulate network issues, or server high load. You must write `delay` as a string with postfix indicating time unit (see [this](https://golang.org/pkg/time/#ParseDuration) 
for more info about an actual format). Also, you can specify minimum and maximum delays using separator ':', the server respond delay will be chosen at random between these values, default value is no delay at all.

### Create an imposter using regex

* [Using regex on the endpoint](#regex-on-the-endpoint)
* [Using regex on the params](#regex-on-the-params)
* [Using regex on the headers](#regex-on-the-headers)

#### Regex on the endpoint

If we want to define a `regex` in our endpoint, we will need using the [gorilla mux](https://stackoverflow.com/questions/28950606/gorilla-mux-regex) nomenclature, that is
`{anyvariablename:{regex}}`.

In the next example, we have an endpoint configured to match with any kind of [ULID Id](https://cran.r-project.org/web/packages/ulid/vignettes/intro-to-ulid.html):

```json
[
  {
    "request": {
      "method": "GET",
      "endpoint": "/gophers/{_id:[\\w]{26}}"
    },
    "response": {
      "status": 200,
      "headers": {
        "Content-Type": "application/json"
      },
      "body": "{\"data\":{\"type\":\"gophers\",\"id\":\"01D8EMQ185CA8PRGE20DKZTGSR\",\"attributes\":{\"name\":\"Zebediah\",\"color\":\"Purples\",\"age\":55}}}"
    }
  }
]
```

#### Regex on the params:

In the case of use a `regex` in our query params, we also will need using the [gorilla mux](https://stackoverflow.com/questions/28950606/gorilla-mux-regex) nomenclature, that is
`{anyvariablename:{regex}}`.

For example, we want an imposter that only do match if we receive an apiKey as param:

```json
[
  {
    "request": {
      "method": "GET",
      "endpoint": "/gophers/{_id:[\\w]{26}}",
      "params": {
        "apiKey": "{_apiKey:[\\w]+}"
      }
    },
    "response": {
      "status": 200,
      "headers": {
        "Content-Type": "application/json"
      },
      "body": "{\"data\":{\"type\":\"gophers\",\"id\":\"01D8EMQ185CA8PRGE20DKZTGSR\",\"attributes\":{\"name\":\"Zebediah\",\"color\":\"Purples\",\"age\":55}}}"
    }
  }
]
```

#### Regex on the headers:

In this case we will not need the `gorilla mux nomenclature` to write our regex.

Then if we want for example do a restrictive `imposter` that need an `Authorization` header to do match, we could do something like
this:

```json
[
  {
    "request": {
      "method": "GET",
      "endpoint": "/gophers/{id:[\\w]{26}}",
      "headers": {
        "Authorization": "\\w+"
      }
    },
    "response": {
      "status": 200,
      "headers": {
        "Content-Type": "application/json"
      },
      "body": "{\"data\":{\"type\":\"gophers\",\"id\":\"01D8EMQ185CA8PRGE20DKZTGSR\",\"attributes\":{\"name\":\"Zebediah\",\"color\":\"Purples\",\"age\":55}}}"
    }
  }
]
```

So as you can see, in the `headers` option you can write directly your value or your regex.

### Create an imposter using JSON Schema

Sometime, we will need to ensure that the request that we will receive is a valid request, for that reason we can
create an `imposter` that only match with a valid [json schema](https://json-schema.org/).

To do that we will need to define our `json schema` first:

`imposters/schemas/create_gopher_request.json`

```json
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

So with this `json schema`, we expect a `request` like this:

```json
{
    "data": {
        "type": "gophers",
        "attributes": {
            "name": "Zebediah",
            "color": "Purples",
            "age": 55
        }
    }
}
```

Then our `imposter` will be configure in the next way:

````json
[
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
        "status": 201,
        "headers": {
            "Content-Type": "application/json"
        }
    }
  }
]
````

The path where the schema is located is relative to where the imposters are.

### Create an imposter with delay

If we want to simulate a problem with the network, or create a more realistic response, we can use the `delay` property.

`delay` property is a range, between two times that use [parser unit of time of Go language](https://golang.org/pkg/time/#ParseDuration).

For example, we can modify our previous POST call to add a `delay`, we want that our call wait between 1 second to 5 second:

````json
[
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
        "status": 201,
        "headers": {
            "Content-Type": "application/json"
        },
        "delay": "1s:5s"
    }
  }
]
````

### Create an Imposter with dynamic responses

Killgrave allow dynamic responses, with that feature we can use one endpoint and obtain different response or even errors.

We need to understand that to do Killgrave need all the `imposters` order from more restrictive to less, because the server
try to do match with each `request` and stop when his find one that have the requirements to be executed.

For example, think on our previous POST example, we want to create a gopher, but the `request` doesn't comply with the
json schema declared.

````json
[
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
        "status": 201,
        "headers": {
            "Content-Type": "application/json"
        }
    }
  },
  {
      "request": {
          "method": "POST",
          "endpoint": "/gophers"
      },
      "response": {
          "status": 400,
          "headers": {
              "Content-Type": "application/json"
          },
          "body": "{\"errors\":\"bad request\"}"
      }
  }
]
````

So as you can see, we have first of all the imposters with the restrictive, `headers` and `json schema`, but
our last `imposter` is a simple `imposter` that it will match with any call via POST to `/gophers`.

## Contributing
[Contributions](CONTRIBUTING.md) are more than welcome, if you are interested please follow our guidelines to help you get started.

## License
MIT License, see [LICENSE](https://github.com/friendsofgo/killgrave/blob/master/LICENSE)
