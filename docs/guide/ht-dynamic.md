---
prev: ./your-first-imposter
next: ./debug-intro
---

# Use dynamic responses

Killgrave allows dynamic responses. Using this feature, Killgrave can return different responses on the same endpoint.

To do this, all imposters need to be sorted from most restrictive to least. Killgrave tries to match the request with 
each of the imposters in sequence, stopping at the first imposter that matches the request.

In the following example, there are defined multiple imposters for the `POST /gophers` endpoint:

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

Now,

1. Let's say an incoming request does not match the JSON schema specified in the first imposter's `schemaFile`.
2. Therefore, Killgrave skips this imposter and tries to match the request against the next configured imposter. 
3. The next configured imposter is much less restrictive, so the request matches and the associated response is returned.