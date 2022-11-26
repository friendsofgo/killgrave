---
prev: ./advanced
next: ./ht-json
---

# Use regular expressions (regex)

## Regex in the endpoint

Killgrave uses the [gorilla/mux](https://github.com/gorilla/mux) regular expression format for endpoint matching.

In the next example, we have configured an endpoint to match with any kind of 
[ULID ID](https://cran.r-project.org/web/packages/ulid/vignettes/intro-to-ulid.html):

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

## Regex in the query parameters

Killgrave uses the [gorilla/mux](https://github.com/gorilla/mux) regular expression format for query parameter matching.

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

## Regex in the headers:

In this case we will not need the `gorilla mux nomenclature` to write our regex.

In the next example, we have configured an imposter that uses regex to match an `Authorization` header.

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