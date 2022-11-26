---
prev: ./ht-json
next: ./ht-dynamic
---

# Use delayed responses

If you want to simulate a problem with the network, or create a more realistic response, you can use the `delay` property.

The `delay` property can take duration in the [Go ParseDuration format](https://golang.org/pkg/time/#ParseDuration). 
The server response will be delayed by the specified duration.

Alternatively, the `delay` property can take a range of two durations, separated by `:`. 
In this case, the server will respond with a random delay within this range.

With the example below, the response would be delayed between 1 and 5 seconds:

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