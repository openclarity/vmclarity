# Plugin SDK for Python

SDK simplifies the development of scanner plugins used in VMClarity.
It provides a set of libraries that the developers can use to
quickly develop new security scanners.

## Usage

```python
TODO(ramizpolic): Provide installation options
```

## Developer notes

- The scanner should be executed in the container
- The scanner should run REST server defined in [Scanner Plugin OpenAPI specs](../../../openapi.yaml)
- Logs should be available on standard output to allow collection by
  other tools that manage the container lifecycle.

All scanner plugins are run as containers and used via REST server interface.
Developers should ship their scanners as container images that run the REST server.

Configuration for the REST server (and container via env vars) are given in [config.py](plugin/scanner/config.py).

Use `plugin.server` module to create and run the REST server.
Default logger is configured on start to write to standard output.
You can use default logger in your code.

## TODO

- Add testing logic to verify that SDK works
