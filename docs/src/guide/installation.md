---
prev: ./concepts
next: ./getting-started
---

# Installation

> ⚠️ Even though Killgrave is a very robust tool and is being used by some companies in production environments, 
> it's still in initial development. Therefore, 'minor' version numbers are used to signify breaking changes and 'patch'
> version numbers are used for non-breaking changes or bugfixing. As soon as v1.0.0 is released, Killgrave will start 
> to use [SemVer](https://semver.org/) as usual.

You can install Killgrave in different ways, but all of them are very simple:

## Go Toolchain

One of them is of course using `go install`, Killgrave is a Go project and therefore can be compiled using the go toolchain:

```sh
$ go install github.com/friendsofgo/killgrave/cmd/killgrave@{version}
```

*Note that `version` must be replaced by the version that you want to install. 
If left unspecified, the `main` branch will be installed.*

## Homebrew

If you are a macOS user, you can install Killgrave using [Homebrew](https://brew.sh/):

```sh
$ brew install friendsofgo/tap/killgrave
```

⚠️ If you are installing via Homebrew, you always get the latest Killgrave version, we hope to fix this soon!

## Docker

Killgrave is also available through [Docker](https://www.docker.com/).

```sh
$ docker run -it --rm -p 3000:3000 -v $PWD/:/home -w /home friendsofgo/killgrave -host 0.0.0.0
```

- `-p 3000:3000` is used to forward the local port `3000` (Killgrave's default port) to container's port `3000`,
otherwise Killgrave won't be reachable from the host.

- `-host 0.0.0.0` is used to change the Killgrave's default host (`localhost`) to allow Killgrave to listen to and 
respond to incoming requests from outside the container, otherwise Killgrave won't be reachable from the host.

## Other

Windows and Linux users can download binaries from the 
[GitHub Releases](https://github.com/friendsofgo/killgrave/releases) page.