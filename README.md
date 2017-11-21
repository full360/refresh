# Refresh

Simple application to reload the configuration of an application through the
"reload" endpoint.

This application will download the updated configuration files from a storage
(s3) into an specified directory and then do an HTTP request to the `app-url`
using the provided method. This application is meant to be used as a side-car to
the Prometheus Server.

## Running

To run the application use the following command:

    ./refresh

There are various arguments that can be passed, use `-h` display

Current endpoints:

- `/health` will allow for application health checks
- `/app/refresh` will trigger a download and refresh of the configured app
- `/metrics` will display Prometheus metrics

## Building and Releasing

To build the project we have set a make task that'll only build Darwin and Linux
binaries for amd64. Remember to Bump the version inside the `Makefile` before
releasing and that's about it. If we there's a need for a different architecture
we can add it.

    make release

## Tests

Running tests can be performed from the default make command or from the test
target

    make test

## Missing

 - [ ] we are missing tests
 - [ ] docs and code documentation
 - [ ] better metrics that include status codes
 - [ ] better logging that include status codes and any other relevant info
