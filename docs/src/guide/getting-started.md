---
prev: ./installation
next: ./your-first-imposter
---

# Getting started

To start Killgrave, you can simply do it by running:

```sh
$ killgrave
```

While you are welcome to provide your own configuration, Killgrave will default to the following [configuration](/config):

> - **imposters path**: `imposters`
> - **host**: `localhost`
> - **port**: `3000`
> - **CORS**: `[]`
> - **proxy**: `none`
> - **watcher**: `false`
> - **debugger**:
>   - **enabled**: `false`
>   - **address**: `localhost:3030`

## Command-line interface

Killgrave is almost fully configurable through the command line, except for CORS, which can only be configured using the
config file. 

You can find the list with all the available flags at [CLI](/cli) section or by running: 

```sh
$ killgrave -h
```

## Configuration file

In case you are looking for a predictable and reproducible configuration, you can use the `-config` command-line flag
to specify the location of a configuration file, which can be written either in JSON or YAML.

You can find the list with all the available settings at [Config Reference](/config) section.

As you can see, you can configure all the options in a very easy way.

> Historically, the options `imposters_path`, `port` and `host` were mandatory when using a configuration file. 
> However, since the last version (`v0.4.1`), they are no longer needed, so you can simply override those options
> if you actually want to. Furthermore, in previous versions (prior to `v0.4.1`), the `imposters_path` option was used to
> refer to the path to the working directory where Killgrave was executed from, but since the last version (`v0.4.1`) 
> this is relative to the location of the config file.


