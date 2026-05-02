# dagger-github-cache-test

A minimal Go [Echo](https://echo.labstack.com/) hello-world web service built and packaged with [Dagger](https://docs.dagger.io/).

The app listens on `:8080` and responds to `GET /` with `Hello, World!`.

## Prerequisites

- [Dagger CLI](https://docs.dagger.io/install) (engine version pinned in `dagger.json`)
- Docker or another supported container runtime for the Dagger engine to talk to

You do **not** need Go installed locally — Dagger runs the build inside a container.

## Dagger functions

List the functions exposed by the module:

```bash
dagger functions
```

| Function    | Returns     | Description                                                       |
| ----------- | ----------- | ----------------------------------------------------------------- |
| `build`     | `File`      | Compile a static Linux binary of the web service from source.     |
| `container` | `Container` | Package the binary into an Alpine image exposing port 8080.       |
| `publish`   | `String`    | Push the container image to a registry and return the digest ref. |
| `serve`     | `Service`   | Run the container as an ephemeral Dagger service on port 8080.    |

All functions take a `--source` directory argument — pass `.` to use the current repo.

### `build` — compile the binary

Build the statically-linked Linux binary and export it to the host:

```bash
dagger call build --source=. export --path=./server
./server   # only works on linux/amd64
```

### `container` — produce a runtime image

Build the Alpine-based runtime container. Export it as an OCI tarball that you can load into Docker:

```bash
dagger call container --source=. export --path=./image.tar
docker load -i image.tar
```

Or inspect it interactively in a Dagger shell:

```bash
dagger call container --source=. terminal
```

### `publish` — push to a registry

Publish to any OCI registry. The example below uses [ttl.sh](https://ttl.sh), an anonymous, ephemeral registry — no auth required:

```bash
dagger call publish --source=. --address=ttl.sh/hello-echo:1h
```

For an authenticated registry, log in on the host first (`docker login ghcr.io`) and Dagger will reuse those credentials:

```bash
dagger call publish --source=. --address=ghcr.io/<you>/hello-echo:latest
```

### `serve` — run the service locally

Run the service and forward the container's port 8080 to your host:

```bash
dagger call serve --source=. up --ports=8080:8080
```

Then in another terminal:

```bash
curl http://localhost:8080/
# Hello, World!
```

Press `Ctrl+C` to stop.

## Project layout

```
.
├── main.go            # Echo hello-world service
├── go.mod / go.sum    # app module (github.com/labstack/echo/v4)
├── dagger.json        # Dagger module config
└── .dagger/
    └── main.go        # Dagger functions: build / container / publish / serve
```

## Updating the Dagger module

After editing `.dagger/main.go`, regenerate the SDK bindings:

```bash
dagger develop
```
