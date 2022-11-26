---
prev: ./your-first-imposter
next: ./ht-regex
---

# Advanced features

Killgrave has some advanced features:

## CORS

If you want to use Killgrave from a frontend application, you should consider configuring
[cross-origin resource sharing (CORS)](https://developer.mozilla.org/docs/Web/HTTP/CORS).

In the [CORS](/config/#cors) section of the file you can find the following options:

- **methods**: `Access-Control-Request-Method`
- **headers**: `Access-Control-Request-Headers`
- **exposed_headers**: `Access-Control-Expose-Headers`
- **origins**: `Access-Control-Allow-Origin`
- **allow_credentials**: `Access-Control-Allow-Credentials`

You can find further details [here](/config/#cors).

## Proxy

You can use Killgrave in proxy mode using the flags `proxy-mode` and `proxy-url` or their equivalent fields in the 
[configuration file](/config). The following three proxy modes are available:

* `none`: ***(Default)*** Killgrave won't behave as a proxy and will only use the configured imposters.
* `missing`: Killgrave will try to match the request with a configured imposter. Otherwise, it will forward the request
to the configured proxy.
* `all`: Killgrave will forward all the incoming requests to the configured proxy.

The `proxy-url` must be the root path of the proxy server. For instance, if we have an API running on 
`http://example.com/things`, the configured `proxy-url` should be `http://example.com`.

## Secure (HTTP over TLS)

Killgrave has a secure mode that lets you expose your mock servers over secure connections by using HTTP over TLS/SSL
(**HTTPS**).

You can use the `secure` setting to enable the secure mode. Disabled by default.

If enabled, the mock server will use TLS options with a dummy certificate, to make it work with the `HTTPS` protocol.

## Watcher

Killgrave has a file watcher that lets you hot reload your mock servers on every imposters change.

You can use the `watcher` setting to enable the secure mode. Disabled by default.

If enabled, the file watcher will be watching changes on the imposters directory so the mock server will be restarted
on every imposters change. 

## Debugger

Killgrave has an interactive mode that lets you debug the behavior of your mock server.

You can use the `debugger` settings to enable the interactive mode. Disabled by default.

You can find further details [here](/guide/debug-intro). 
