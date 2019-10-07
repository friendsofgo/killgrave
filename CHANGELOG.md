# Changelog

## v0.4.0 (TBD)
* The config file option load the imposters path relative on where the config file is
* Upgrade Killgrave to go1.13

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
