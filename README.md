# license-server

Server and CLI for managing software licenses.

## CLI Installation

**Note**: The CLI uses the `getConfig` function in `./internal/config/config.go` which loads the application configuration file. The license manager CLI must be used from a directory containing the file or a subdirectory.

```
make cli
licmgr --version
```
