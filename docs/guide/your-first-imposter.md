---
prev: ./getting-started
next: ./advanced
---

# Your first imposter

**At least one imposter must be configured** in order to use Killgrave.

Every file with any of the valid extensions (`.imp.json`, `.imp.yml` or `.imp.yaml`), present in the "imposters" folder
(default `"./imposters"`) will be interpreted as an imposters file.

We use a rule-based system to match requests to imposters. Therefore, you have to organize your imposters **from most
restrictive to least**.

Here's an example of the contents of an imposters file with a single imposter:

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

This a very simple example. Killgrave has more possibilities for configuring imposters. 

You can take a look at some of them in the **"How to...?"** section below.

> ⚠️Remember that you will need to escape any special char, in the properties that admit text.