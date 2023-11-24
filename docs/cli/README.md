---
sidebar: auto
---

# Command-line interface (CLI)

<img :src="$withBase('/img/killgrave.png')" alt="killgrave" style="max-width:130px">

## Killgrave

Killgrave is basically a command-line interface (CLI) that can be used with no explicit configuration, but a set
of [imposters](/config/#imposters). Look at the [config reference](/config) or [guide](/guide) for further details.

However, you can tune up some of their settings like the host and port where the mock server is listening to, 
among others, by providing some configuration settings.

To provide those settings, you can either use the [available CLI flags](#available-flags) or use the `-config` 
flag to provide the path to a settings file. In such case, you can either use a JSON or YAML configuration file.

### Available flags

See below the list of available flags:

```sh
$ killgrave -h

  -config string
        path to the configuration file
  -debugger
        run your server with the debugger
  -debugger-addr string
        debugger address (default "localhost:3030")
  -host string
        run your server on a different host (default "localhost")
  -imposters string
        directory where imposters are read from (default "imposters")
  -port int
        run your server on a different port (default 3000)
  -proxy-mode string
        proxy mode (choose between 'all', 'missing' or 'none') (default "none")
  -proxy-url string
        proxy url, use it in combination with proxy-mode
  -secure
        run your server using TLS (https)
  -version
        show the version of the application
  -watcher
        enable the file watcher, which reloads the server on every file change
```