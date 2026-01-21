# BubbleWebServer

Go TUI + simple HTTP server packaged as an importable module.

This repository exposes a `bubbleserver` package that provides a BubbleTea-based TUI to manage a small HTTP server and register routes.

## Quick start

1. Ensure Go is installed (Go 1.25+ recommended).

   ```bash
   # install the bubbleserver package
   go get github.com/efraim132/BubbleWebServer/bubbleserver
   ```

3. See `example_main.go` for a full usage example. Example imports:

   ```go
   import "github.com/efraim132/BubbleWebServer/bubbleserver"
   ```

## Running locally

- For quick testing run `go run main.go` from the repo root. This runs the local TUI + bundled webserver.
- Default webserver port is 8090; use `/set port` in the TUI or call `server.SetDefaultPort(...)` in code.

## Screenshots

TUI main screen:

![TUI Main Screen](assets/tui-main-screen.png)

Port selection form:

![TUI Port Selection](assets/tui-port-selection.png)

