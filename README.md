# Prometheus Reload

Simple application to reload Prometheus on configuration changes. This
application is meant to be used as a side-cart to the Prometheus Server.

## Requirements

- HTTP API
  - Endpoint for health checks
  - Endpoint for reload notifications

## Flags

- Address (addr)
- Port (port)
- Service Name (svcName)
- Prometheus Server Address (promAddr)
- Endpoint to trigger the reload (endpoint)
- AWS S3 region (awsRegion)
- AWS S3 Bucket name (s3Bucket)
- Download directory (dlDir)
