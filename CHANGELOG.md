# Changelog

## v0.5.0 (2023/05/12)
* Replace the use of `go flags` by `cobra`
* Change all the `cli` commands and flags, more information on `README.md` (Breaking change if you are using single dash)
* Deprecated Go Version 1.16, we support now Go 1.20
* Change support of the filesystem, start using [afero](https://github.com/spf13/afero) library
* Deprecated using test on the Go standard way and start to use [testify](https://github.com/stretchr/testify) library
* Fix a bug with the `proxy-mode` `missing`, the app always tryed to load the imposter even when not match
* Using stric slashes

## v0.4.1(2021/04/24)
* Migration to Github actions and remove the use of Circle CI
* Deprecation go versions before to v1.16
* Run the mock server using TLS
* Support YAML for a configuration and imposters definition
* Put watcher options on the config file
* Not support codecov anymore

## v0.4.0 (2020/01/30)
* The config file option load the imposters path relative on where the config file is
* Upgrade Killgrave to go1.13
* Remove use of github.com/pkg/errors in favor to standard errors package
* Remove backward compatibility with previous versions to go 1.13
* Add `-watcher` flag to reload the server with any changes on `imposters` folder
* Fix searching imposter files mechanism
* Add proxy server feature
* Allow to add latency to responses

## v0.3.3 (2019/05/11)

* Improve default CORS options
* Allow up mock server via config file
* Allow configure CORS options
  * Access-Control-Request-Method
  * Access-Control-Request-Headers
  * Access-Control-Allow-Origin
  * Access-Control-Expose-Headers
  * Access-Control-Allow-Credentials
* Improve route_mateches unit tests

## v0.3.2 (2019/05/08)

* Fix CORS add AccessControl allowing methods and headers

## v0.3.1 (2019/05/07)

* Allow CORS requests

## v0.3.0 (2019/04/28)

* Dynamic responses based on headers
* Standarize json files using [Google JSON style Guide](https://google.github.io/styleguide/jsoncstyleguide.xml)
* Move to `internal` not exposable API
* Dynamic responses based on query params
* Allow organize your imposters with structured folders (using new extension `.imp.json`)
* Allow write multiple imposters by file

## v0.2.1 (2019/04/25)

* Allow imposter's matching by request schema
* Dynamic responses based on regex endpoint or request schema
* Calculate files directory(body and schema) based on imposters path
* Update REAMDE.md with resolved features and new future features

## v0.2.0 (2019/04/24)

* Create an official docker image for the application
* Update README.md with how to use the application with docker
* Allow write headers for the response

## v0.1.0 (2019/04/21)

* Add Killgrave logo
* Add CircleCI integration
* Convert headers into canonical mime type
* Run server with imposter configuration
* Processing and parsing imposters file
* Initial version
