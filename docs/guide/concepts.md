---
prev: ./
next: ./installation
---

# Concepts

## Imposters 

Imposters are the most important concept in the Killgrave's world.

They conform the rules that determine how the mock server should respond to a request.

You can identify a Killgrave imposter file by its extension: `.imp.json`, `.imp.yml` or `.imp.yaml`.

> *You can learn more about how to configure imposters in the [Imposter Configuration Section](/config/#imposters).*

### Imposters Structure

The imposter object can be divided in two parts:

- [Request](#request)
- [Response](#response)

#### Request

This part defines how Killgrave should determine whether an incoming request matches the imposter or not. 

The request object has the following properties:

- `method` *(mandatory)*: The HTTP method of the incoming request.
- `endpoint` *(mandatory)*: Path of the endpoint relative to the base. Supports regular expressions.
- `schemaFile`: A JSON schema to validate the incoming request against.
- `params`: Restrict incoming requests by query parameters. Supports regular expressions.
- `headers`: Restrict incoming requests by HTTP header.

#### Response

This part defines how Killgrave should, in case of match, respond to the incoming request.

The response object has the following properties:

- `status` *(mandatory)*: Integer defining the HTTP status code to return.
- `body` or `bodyFile`: The response body. Either a literal string (`body`) or a path to a file (`bodyFile`). `bodyFile`
is especially useful in the case of large outputs. This property is optional: if not response body should be returned it
should be removed or left empty.
- `headers`: HTTP headers to return in the response.
- `delay`: Time the server waits before responding. This can help simulate network issues, or high server load. 
Uses the Go [`time.ParseDuration`](https://pkg.go.dev/time#ParseDuration) format. Also, you can specify minimum and 
maximum delays separated by  `:`. The response delay will be chosen randomly between these values. 
Default value is `"0s"` (no delay).