---
prev: ./ht-regex
next: ./ht-delays
---

# Use JSON Schema

Sometimes you may need to validate requests more thoroughly. In such case, you can
create an imposter that only matches with a valid [json schema](https://json-schema.org/).

To do that you will need to define our `json schema` first:

*e.g. `imposters/schemas/create_gopher_request.json`*

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

Then the imposter could be configured as follows:

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