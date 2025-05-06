# Biathlon

## Description
This is a CLI app for parsing biathlon competition logs. It parses events from a source file and prints results of a competition in console terminal.
## Run

To run app, you can use Makefile command:
```bash
make test_run
```
The output will appear in output.txt file.

If you want to pass your own files, you can use usual go tools:

```bash
go run ./cmd/main.go 'input events file', 'config.json file', 'destination file' 
```

## Test
```bash
go test -v ./internal/...
```
