---
sidebar: auto
---

# Config Reference

<img :src="$withBase('/img/killgrave.png')" alt="killgrave" style="max-width:130px">

## Killgrave

Killgrave can be used without explicitly providing any configuration. However, you can tune up some of their settings
like the host and port where the mock server is listening to, among others, by providing some configuration settings.

To provide those settings, you can either use the available [CLI](/cli) flags or use the `-config` flag to provide the
path to a settings file. In such case, you can either use JSON or YAML.

### Basic

Minimal configuration attributes needed to get your mock up and running.

#### host

- Type: `string`
- Default: `localhost`

Specify the host for the Killgrave server. If you are using [Docker](https://www.docker.com/), you must override
the default value. Otherwise, the server will be reachable only from the container itself, so probably it won't work
as you expect.

```json
{
  "host": "localhost"
}
```

```yaml
host: "localhost"
```

#### port

- Type: `number`
- Default: `3000`

Specify the port for the Killgrave server. If you are using [Docker](https://www.docker.com/), you need to forward it.

```json
{
  "port": 3000 
}
```

```yaml
port: 3000
```

#### imposters_path

- Type: `string`
- Default: `imposters`

Specify the directory the imposter files (either `.imp.json`, `.imp.yml` or `.imp.yaml`) will be loaded from. 
On a regular set up, this directory will contain multiple of those imposter files and the directories with 
`schemas` and `responses`.

```json
{
  "imposters_path": "imposters" 
}
```

```yaml
imposters_path: imposters
```

### Proxy

Set up a proxy to redirect the incoming requests as you want.

#### mode

- Type: `string`
- Default: `none`

Specify the proxy mode for the Killgrave server. The default value is `none` which means no proxy.
Use `missing` to redirect only those incoming requests that aren't defined within the imposters. Use `all` to redirect all incoming requests.

```json
{
  "proxy": {
    "mode": "missing"
  }
}
```

```yaml
proxy:
  mode: missing
```

#### url

- Type: `string`
- Default: -

Specify the url for the Killgrave's proxy. The incoming requests will be redirected to this url based on the proxy mode.

```json
{
  "proxy": {
    "url": "https://example.com"
  }
}
```

```yaml
proxy:
  url: https://example.com
```

### CORS

Set up the [cross-origin resource sharing (CORS)](https://developer.mozilla.org/docs/Web/HTTP/CORS) mechanism for the 
Killgrave server. Especially useful when mocking servers that are consumed by frontend applications. 

#### methods

- Type: `string array`
- Default: `["GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE", "PATCH", "TRACE", "CONNECT"]`

Represents the `Access-Control-Request-Method` header.

```json
{
  "cors": {
    "methods": ["GET"]
  }
}
```

```yaml
cors:
  methods: ["GET"]
```

#### headers

- Type: `string array`
- Default: `["X-Requested-With", "Content-Type", "Authorization"]`

Represents the `Access-Control-Request-Headers` header.

```json
{
  "cors": {
    "headers": ["Content-Type"]
  }
}
```

```yaml
cors:
  headers: ["Content-Type"]
```

#### exposed_headers

- Type: `string array`
- Default: `["Cache-Control", "Content-Language", "Content-Type", "Expires", "Last-Modified", "Pragma"]`

Represents the `Access-Control-Expose-Headers` header.

```json
{
  "cors": {
    "exposed_headers": ["Cache-Control"]
  }
}
```

```yaml
cors:
  exposed_headers: ["Cache-Control"]
```

#### origins

- Type: `string array`
- Default: `[]`

Represents the `Access-Control-Allow-Origin` header.

```json
{
  "cors": {
    "origins": ["*"]
  }
}
```

```yaml
cors:
  origins: ["*"]
```

#### allow_credentials

- Type: `boolean`
- Default: `false`

Enables or disables the `Access-Control-Allow-Credentials` header.

```json
{
  "cors": {
    "allow_credentials": true
  }
}
```

```yaml
cors:
  allow_credentials: true
```

### Security

Sometimes you may want to simulate a production-like behavior with your mock servers, for instance to verify that your
frontend application manages correctly HTTPS/TLS connections.

#### secure

- Type: `boolean`
- Default: `false`

Expose the mock server through HTTP over TLS(SSL).

```json
{
  "secure": true
}
```

```yaml
secure: true
```

### Hot reloads

When building the imposters to mock your server, you way want to reduce the feedback loop on configuration changes.

#### watcher

- Type: `boolean`
- Default: `false`

Enable the file watcher to hot reload the mock server on every imposters change.

```json
{
  "watcher": true
}
```

```yaml
watcher: true
```

### Full example

See below a full example of the configuration file:

```json
{
    "port": 3000,
    "host": "localhost",
    "imposters_path": "imposters",
    "proxy": {
      "url": "https://example.com",
      "mode": "missing"
    },
    "cors": {
      "methods": [
        "GET"
      ],
      "headers": [
        "Content-Type"
      ],
      "exposed_headers": [
        "Cache-Control"
      ],
      "origins": [
        "*"
      ],
      "allow_credentials": true
    },
    "secure": true,
    "watcher": true
}
```

```yaml
port: 3000
host: "localhost"
imposters_path: "imposters"
proxy:
  url: https://example.com
  mode: missing
cors:
  methods: ["GET"]
  headers: ["Content-Type"]
  exposed_headers: ["Cache-Control"]
  origins: ["*"]
  allow_credentials: true
secure: true
watcher: true
```

## Imposters

Imposters are the first-class citizens in Killgrave, they define how the mock server will respond to the incoming requests.
You can use either JSON or YAML to define them. 

### Request

The `request` configuration is used to evaluate whether an incoming request matches the imposter or not.
If there's a match, the mock server will respond with the imposter's response.

#### method

- Type: `string`
- Default: -

Specify the expected HTTP method.

```json
{
  "request": {
    "method": "POST"
  }
}
```

```yaml
request:
  method: "POST"
```

#### endpoint

- Type: `string`
- Default: -

Specify the expected URL endpoint (path). You can define URL parameters and use regular expressions.
Look at [Gorilla Mux examples](https://github.com/gorilla/mux#examples) to see how regular expressions work.

```json
{
  "request": {
    "endpoint": "/gophers"
  }
}
```

```yaml
request:
  endpoint: /gophers
```

#### schemaFile

- Type: `string`
- Default: -

Specify the path to the file that contains the expected JSON schema.

```json
{
  "request": {
    "schemaFile": "schemas/create_gopher_request.json"
  }
}
```

```yaml
request:
  schemaFile: "schemas/create_gopher_request.json"
```

#### params

- Type: `string map`
- Default: -

Specify the expected URL params. You can use regular expressions.
Look at [Gorilla Mux examples](https://github.com/gorilla/mux#examples) to see how regular expressions work.

```json
{
  "request": {
    "params": {
      "id": "01EKPT"
    }
  }
}
```

```yaml
request:
  params:
    id: "01EKPT"
```

#### headers

- Type: `string map`
- Default: -

Specify the expected HTTP headers.

```json
{
  "request": {
    "headers": {
      "Content-Type": "application/json",
      "Return-Error": "error"
    }
  }
}
```

```yaml
request:
  headers:
    Content-Type: "application/json"
    Return-Error: "error"
```

### Response

The `response` configuration is used to build the response in case of match with the incoming request.

#### status

- Type: `number`
- Default: `200`

Specify the response's HTTP status code.

```json
{
  "response": {
    "status": 401
  }
}
```

```yaml
response:
  status: 401
```

#### body

- Type: `string`
- Default: -

Specify the raw contents (as string) to be returned as the response body. You could use a stringified JSON, for instance.

Although, in such cases where the desired response body is a complex payload (e.g. a JSON object), we recommend to use a
separate file in combination with the [`bodyFile`](#bodyfile) request attribute. 

```json
{
  "response": {
    "body": "Simple response body"
  }
}
```

```yaml
response:
  body: "Simple response body"
```

#### bodyFile

- Type: `string`
- Default: -

Specify the path to the file that contains the response body's contents.

```json
{
  "response": {
    "bodyFile": "/path/to/file"
  }
}
```

```yaml
response:
  bodyFile: "/path/to/file"
```

#### headers

- Type: `string map`
- Default: -

Specify the response's HTTP headers.

```json
{
  "response": {
    "headers": {
      "Content-Type": "application/json",
      "Return-Error": "error"
    }
  }
}
```

```yaml
response:
  headers:
    Content-Type: "application/json"
    Return-Error: "error"
```

#### delay

- Type: `string`
- Default: -

Specify the response delay. It is really helpful to reproduce real-world examples and/or to simulate situations with
peaks of load where the responses are considerably slow.  

You can use either a fixed time (single value) or an interval (two values joined by `:`).

Look at Go [`time.ParseDuration`](https://pkg.go.dev/time#ParseDuration) examples to see the accepted
format to specify the time durations of the response delay. 

```json
{
  "response": {
    "delay": "2s:5s"
  }
}
```

```yaml
response:
  delay: "2s:5s"
```

### Full examples

See below some full examples:

#### Request

```json
{
    "request": {
      "method": "POST",
      "endpoint": "/gophers",
      "schemaFile": "schemas/create_gopher_request.json",
      "params": {
        "id": "01EKPT"
      },
      "headers": {
        "Content-Type": "application/json",
        "Return-Error": "error"
      }
    }
}
```

```yaml
request:
  method: "POST"
  endpoint: "/gophers"
  schemaFile": "schemas/create_gopher_request.json"
  params":
    id: "01EKPT"
  headers:
    Content-Type: "application/json"
    Return-Error: "error"
```

#### Response

```json
{
    "response": {
      "status": 401,
      "body": "Simple response body",
      "bodyFile": "/path/to/file",
      "headers": {
        "Content-Type": "application/json",
        "Return-Error": "error"
      },
      "delay": "2s:5s"
    }
}
```

```yaml
response:
  status: 401
  body: "Simple response body"
  bodyFile: "/path/to/file"
  headers:
    Content-Type: "application/json"
    Return-Error: "error"
  delay: "2s:5s"
```

*Note from the examples above that you should only provide one of [`body`](#body) or [`bodyFile`](#bodyfile), but not
both. Here both are defined just for the sake of showing a full example.*