# mb-tui

A terminal UI for running a handful of internal deployment and release scripts,
built on [Bubble Tea](https://github.com/charmbracelet/bubbletea).

Launching the app presents a searchable menu of scripts:

- **View All Versions** — fetches deployed versions over HTTP
- **Wwwinc Setup Integration** — batches and copies files to the wwwinc servers over SFTP
- **Peak Setup Integration** — the same flow for the Peak servers

## Requirements

- Go 1.22 or newer (see `go.mod`)

## Building

The `main` package lives in `app/`, not at the repo root, so plain `go build .`
will not work — there are no Go files at the root for it to compile. Name the
package explicitly instead:

```bash
go build -o ti ./app
```

That produces a `ti` binary in the repo root, which is already listed in
`.gitignore`. Without `-o`, Go names the binary after its directory and you
would get one called `app`.

To build and run in a single step, without leaving a binary behind:

```bash
go run ./app
```

## Checking the whole module

`go build ./app` only compiles the packages that `app` actually imports. To
type-check every package in the module, including any that are temporarily
unreferenced:

```bash
go build ./...
go vet ./...
```

This is worth doing before committing. Entries in `scripts/` are sometimes
commented out of `app/main.go` while work is in progress, and errors in those
packages will not surface from an `./app` build alone.

## Layout

```
app/         package main — entrypoint and the top-level model that
             switches between scripts
tui/         shared UI components (searchable menu, shell command runner)
tui/downloader/
             download progress UI
exec/        command execution helpers
scripts/     one package per script, each exposing a Bubble Tea model
```

A note on naming: the files at `scripts/*/main.go` are **not** entrypoints.
Despite the filename they declare library packages such as
`package peak_setup_integration`, and the only `package main` in the module is
in `app/`.

## Keybindings

The menu uses readline-style navigation:

| Key | Action |
| --- | --- |
| `ctrl+n` / `ctrl+p` | move down / up |
| `ctrl+g` | clear the search filter |
| `enter` | select |
| `esc` / `ctrl+c` | quit |
