# appearance-watcher

Although most CLI tools are able to automatically switch between light & dark modes, some only statically define their theming.
This tool is aimed at giving automation to such statically defined tools by watching for appearance changes and do such actions like writing to a file.

## Building

This tool uses [goyek](https://github.com/goyek/goyek) to test, lint and build the binary (saved to `./bin`).

To run tests and build the binary, run:

```shell
go run ./build
```

To just build the binary, run:

```shell
go run ./build build
```

## Example configuration

```yaml
symlinks:
  - path: "./file.txt"
    light: "./light.txt"
    dark: "./dark.txt"

files:
  - path: "./content.txt"
    light: "light"
    dark: "dark"
```
