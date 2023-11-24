---
prev: false
next: ./concepts.md
---

# Killgrave

<img :src="$withBase('/img/killgrave.png')" alt="killgrave" style="max-width:130px">

**Killgrave** is a tool that provides a simple way to create a powerful simulator for HTTP-based APIs.

The Killgrave's philosophy is to provide an easy way to configure your mock server, ensuring that you don't need to 
learn another tooling language. For that reason we use JSON and YAML to generate all necessary configurations.

Some Killgrave features are:

- An easy way to create imposters files, either using JSON or YAML.
- The possibility to validate requests against JSON schemas.
- Use of regular expressions to allow different parameters in headers and urls.
- Custom and dynamic response bodies, of all content types (application/json, text/html...).
- Custom response headers and cross-origin resource sharing (CORS).
- Simulate network issues and server high loads with configurable, dynamic responses delay.
- Use of configurable proxy server to call to the original server.
- Use either command-line interface flags or a configuration file (JSON/YAML).
- Hot configuration reloads based on a configurable file watcher.
- Expose the mocker server through HTTP over TLS (SSL).
- Use of an interactive mode to debug your mock server with a web app.