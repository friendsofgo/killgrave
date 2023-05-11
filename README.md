<p align="center">
  <img src="https://res.cloudinary.com/fogo/image/upload/c_scale,w_350/v1555701634/fogo/projects/gopher-killgrave.png" alt="Golang Killgrave"/>
</p>

# Killgrave

Killgrave is a simulator for HTTP-based APIs, in simple words a **Mock Server**, very easy to use, made in **Go**.

![Github actions](https://github.com/friendsofgo/killgrave/actions/workflows/main.yaml/badge.svg?branch=main)
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
    * [Homebrew](#homebrew)
    * [Docker](#docker)
    * [Compile by yourself](#compile-by-yourself)
    * [Other](#other)
- [Getting Started](#getting-started)
    * [Using Killgrave by command line](#using-killgrave-from-the-command-line)
    * [Using Killgrave by config file](#using-killgrave-by-config-file)
    * [Configure CORS](#configure-cors)
    * [Preparing Killgrave for Proxy Mode](#preparing-killgrave-for-proxy-mode)
    * [Creating an Imposter](#creating-an-imposter)
    * [Imposters structure](#imposters-structure)
    * [Using regex in imposters](#using-regex-in-imposters)
    * [Creating an imposter using JSON Schema](#creating-an-imposter-using-json-schema)
    * [Creating an imposter with delay](#creating-an-imposter-with-delay)
    * [Creating an imposter with dynamic responses](#creating-an-imposter-with-dynamic-responses)
- [Contributing](#contributing)
- [License](#license)

## Overview

**Killgrave** is a tool that provides a simple way to create a powerful simulator for HTTP-based APIs.

The Killgrave's philosophy is to provide an easy way to configure your *mock server*, ensuring that you don't need to learn
another tooling language. For that reason we use `json` and `yaml` to generate all necessary configurations.

Killgrave provides:

* An easy way to create imposters files, using `json`
* The possibility to validate requests against [json schemas](https://json-schema.org/).
* Validation of request headers.
* Using regex to allow different parameters in headers and urls.
* Custom body and dynamic body responses.
* Using all content-types bodies, (`application/json`, `text/html`, `text/plain`, `image/jpeg`, etc. )
* Configure response headers.
* Configure CORS.
* Simulate network issues and server high loads with imposter responses delay.
* Using configurable proxy server to call to the original server.
* Run the tool using flags or using a config file.
* Run your mock server using a watcher to reload on configuration changes.

## Concepts

### Imposters

Imposters are the most important concept of the Killgrave tool. They define the rules that determine how the server should respond to a request.

You can identify a Killgrave imposter file by its extension: `.imp.json`.

You can learn more about how to configure imposters in the [Imposter Configuration Section](#imposters).

## Installing
> :warning:  Even though Killgrave is a very robust tool and is being used by some companies in production environments, it's still in initial development. Therefore, 'minor' version numbers are used to signify breaking changes and 'patch' version numbers are used for non-breaking changes or bugfixing. As soon as v1.0.0 is released, Killgrave will start to use [`SemVer`](https://semver.org/) as usual.

You can install Killgrave in different ways, but all of them are very simple:

### Homebrew 

If you are a Mac user, you can install Killgrave using [Homebrew](https://brew.sh/):

```sh
$ brew install friendsofgo/tap/killgrave
```

:warning:  If you are installing via Homebrew, you always get the [latest Killgrave version](https://github.com/friendsofgo/killgrave/releases), we hope to fix this soon.

### Docker

The application is also available through [Docker](https://hub.docker.com/r/friendsofgo/killgrave).

```bash
docker run -it --rm -p 3000:3000 -v $PWD/:/home -w /home friendsofgo/killgrave -host 0.0.0.0
```

`-p 3000:3000` [publishes](https://docs.docker.com/engine/reference/run/#expose-incoming-ports) port 3000 (Killgrave's default port) inside the
container to port 3000 on the host machine.

`-host 0.0.0.0` is necessary to allow Killgrave to listen and respond to requests from outside the container (the default,
`localhost`, will not capture requests from the host network).

### Compile by yourself

If you want to use `Killgrave` from the source code, first you will need to clone the repository:

```sh
git clone git@github.com:friendsofgo/killgrave.git
```

Select the branch you want to use (`main` by default), and then compile `Killgrave`:

```sh
make build
```

This command will create an executable into the path `bin/killgrave`. We highly recommended the usage of this command
to compile the software because parameters such as `version` are added to the compilation, which are necessary when reporting a bug.

### Other

Windows and Linux users can download binaries from the [Github Releases](https://github.com/friendsofgo/killgrave/releases) page.

## Getting Started

To start Killgrave, you simply run the following.

```sh
$ killgrave
```

While you are welcome to provide your own configuration, Killgrave will default to the following configuration:

* **imposters path**: `imposters`
* **host**: `localhost`
* **port**: `3000`
* **CORS**: `[]`
* **proxy**: `none`
* **watcher**: `false`

### Using Killgrave from the command line

Killgrave takes the following command line options. Killgrave is almost fully configurable through the command line, except for `CORS`, which can only be configured using the [config file](#using-killgrave-by-config-file).

```sh
$ killgrave -h

Simple way to generate mock servers

Usage:
  killgrave [flags]

Flags:
  -c, --config string       Path to your configuration file
  -h, --help                Help for Killgrave
  -H, --host string         Set a different host than localhost (default "localhost")
  -i, --imposters string    Directory where your imposters are located (default "imposters")
  -P, --port int            Port to run the server (default 3000)
  -m, --proxy-mode string   Proxy mode, the options are all, missing or none (default "none")
  -u, --proxy-url string    The url where the proxy will redirect to
  -s, --secure              Run mock server using TLS (https)
  -v, --version             Version of Killgrave
  -w, --watcher             File watcher will reload the server on each file change
```

### Using Killgrave by config file

If we want a more permanent configuration, we could use the option `-config` to specify the location of a configuration file.

The config file must be a YAML file with the following structure.

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

As you can see, you can configure all the options in a very easy way. For the above example, the file tree looks as follows, with the current working directory being `mymock`.

```
mymock/
    imposters/
        config.yml
        swapi_people.imp.json
        swapi_planets.imp.json
    Dockerfile
    MakeFile
```

Then in your config file, you will need to set the `-imposters` flag to `.` because the imposters folder is located in the same folder.

Historically, the options `imposters_path`, `port`, `host` were mandatory when using a configuration file. However, since the last version, they are no longer needed, so you can simply override those options if you want to.
Furthermore, in previous versions, the `imposters_path` option refered to the path where the app was launched, but in the last version this is relative to the location of the config file.

The option `cors` is still optional and its options can be an empty array.
If you want more information about the CORS options, visit the [CORS section](#configure-cors).

The `watcher` configuration field is optional. With this setting you can enable hot-reloads on imposter changes. Disabled by default.

The `secure` configuration field is optional. With this setting you can run your server using TLS options with a dummy certificate, so as to make it work with the `HTTPS` protocol. Disabled by default.

The option `proxy-mode` allows you to configure the mock in proxy mode. When this mode is enabled, Killgrave will forward any unconfigured requests to another server. More information: [Proxy Section](#prepare-killgrave-for-proxy-mode)

## How to use

### Configure CORS

If you want to use `killgrave` from a client application you should consider configuring CORS.

In the CORS section of the file you can find the following options:

- **methods** (string array)
  
  Represents the **Access-Control-Request-Method header**, if you don't specify it or if you leave it as an empty array, the default values will be:

  `"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE", "PATCH", "TRACE", "CONNECT"`

- **headers** (string array)
  
  Represents the **Access-Control-Request-Headers header**, if you don't specify it or if you leave it as an empty array, the default values will be:

  `"X-Requested-With", "Content-Type", "Authorization"`

- **exposed_headers** (string array)
  
  Represents the **Access-Control-Expose-Headers header**, if you don't specify it or if you leave it as an empty array, the default values will be:

  `"Cache-Control", "Content-Language", "Content-Type", "Expires", "Last-Modified", "Pragma"`

- **origins** (string array)
  
  Represents the **Access-Control-Allow-Origin header**, if you don't specify it or if you leave it as an empty array this option has not default value

- **allow_credentials** (boolean)
  
  Enables or disables the **Access-Control-Allow-Credentials header**.

### Preparing Killgrave for Proxy Mode

You can use Killgrave in proxy mode using the flags `proxy-mode` and `proxy-url` or their equivalent fields in the configuration file. The following three proxy modes are available:
* `none`: Default. Killgrave will not behave as a proxy and the mock server will only use the configured imposters.
* `missing`: With this mode the mock server will try to match the request with a configured imposter, but if no matching endpoint was found, the mock server will call to the real server, declared in the `proxy-url` configuration variable.
* `all`: The mock server will always call to the real server, declared in the `proxy-url` configuration variable.

The `proxy-url` must be the root path of the proxied server. For example, if we have an API running on `http://example.com/things`, the `proxy-url` will be `http://example.com`.

### Creating an Imposter

At least one imposter must be configured in order to run Killgrave. Files with the `.imp.json` extension in the `imposters` folder (default "imposters") will be interpreted as imposter files.

We use a rule-based system to match requests to imposters. Therefore, you have to organize your imposters from most restrictive to least. Here's an example of an imposter.

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

This a very simple example. Killgrave has more possibilities for configuring imposters. Let's take a look at some of them in the next sections.

:warning:  Remember that you will need to escape any special char, in the properties that admit text.

### Imposters Structure

The imposter object can be divided in two parts:

* [Request](#request)
* [Response](#response)

#### Request

This part defines how Killgrave should determine whether an incoming request matches the imposter or not. The `request` object has the following properties:

* `method` (<span style="color:red">mandatory</span>): The [HTTP method](https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods) of the incoming request.
* `endpoint` (<span style="color:red">mandatory</span>): Path of the endpoint relative to the base. Supports regex.
* `schemaFile`: A JSON schema to validate the incoming request against.
* `params`: Restrict incoming requests by query parameters. More info can be found [here](#create-an-imposter-with-query-params). Supports regex.
* `headers`: Restrict incoming requests by HTTP header. More info can be found [here](#create-an-imposter-with-headers).

#### Response

This part defines how Killgrave should respond to the incoming request. The `response` object has the following properties:

* `status` (<span style="color:red">mandatory</span>): Integer defining the HTTP status to return.
* `body` or `bodyFile`: The response body. Either a literal string (`body`) or a path to a file (`bodyFile`). `bodyFile` is especially useful in the case of large outputs.
This property is optional: if not response body should be returned it should be removed or left empty.
* `headers`: Headers to return in the response.
* `delay`: Time the server waits before responding. This can help simulate network issues, or high server load. Uses the [Go ParseDuration format](https://pkg.go.dev/time#ParseDuration). Also, you can specify minimum and maximum delays separated by ':'. The response delay will be chosen at random between these values. Default value is "0s" (no delay).

### Using regex in imposters

* [Using regex in the endpoint](#regex-on-the-endpoint)
* [Using regex in the query parameters](#regex-on-the-params)
* [Using regex in the headers](#regex-on-the-headers)

#### Regex in the endpoint

Killgrave uses the [gorilla/mux](https://github.com/gorilla/mux) regex format for endpoint regex matching.

In the next example, we have configured an endpoint to match with any kind of [ULID ID](https://cran.r-project.org/web/packages/ulid/vignettes/intro-to-ulid.html):

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

#### Regex in the query parameters:

Killgrave uses the [gorilla/mux](https://github.com/gorilla/mux) regex format for query parameter regex matching.

In this example, we have configured an imposter that only matches if we receive an apiKey as query parameter:

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

#### Regex in the headers:

In this case we will not need the `gorilla mux nomenclature` to write our regex.

In the next example, we have configured an imposter that uses regex to match an Authorization header.

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

### Creating an imposter using JSON Schema

Sometimes, we need to validate our request more thoroughly. In cases like this we can
create an imposter that only matches with a valid [json schema](https://json-schema.org/).

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

With this `json schema`, we expect a `request` like this:

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

Then our imposter will be configured as follows:

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

### Creating an imposter with delay

If we want to simulate a problem with the network, or create a more realistic response, we can use the `delay` property.

The `delay` property can take duration in the [Go ParseDuration format](https://golang.org/pkg/time/#ParseDuration). The server response will be delayed by the specified duration.

Alternatively, the `delay` property can take a range of two durations, separated by a ':'. In this case, the server will respond with a random delay in this range.

For example, we can modify our previous POST call to add a `delay` to determine that we want our response to be delayed by 1 to 5 seconds:

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

### Creating an imposter with dynamic responses

Killgrave allows dynamic responses. Using this feature, Killgrave can return different responses on the same endpoint.

To do this, all imposters need to be ordered from most restrictive to least. Killgrave tries to match the request with each of the imposters in sequence, stopping at the first imposter that matches the request.

In the following example, we have defined multiple imposters for the `POST /gophers` endpoint. Let's say an incoming request does not match the JSON schema specified in the first imposter's `schemaFile`. Therefore, Killgrave skips this imposter and tries to match the request against the second imposter. This imposter is much less restrictive, so the request matches and the associated response is returned.

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

## Contributing
[Contributions](CONTRIBUTING.md) are more than welcome, if you are interested please follow our guidelines to help you get started.

## License
MIT License, see [LICENSE](https://github.com/friendsofgo/killgrave/blob/main/LICENSE)
