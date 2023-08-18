# Build:
```
GOARCH=amd64 GOOS=windows GO111MODULE=on CGO_ENABLED=0 go build -o ./processlist-OTE-0.79.1-exclude-thread-fix.exe
```

# Usage:
## Help:
```
.\processlist-OTE-0.79.1-exclude-thread-fix.exe -h
Usage of C:\Users\Administrator\Downloads\processlist-OTE-0.79.1-exclude-thread-fix.exe:
  -log string
        debug : Display additional output (default "info")
  -version
        Display version

```
## Version:
To view the module versions, run it with --version
```
.\processlist-OTE-0.79.1-exclude-thread-fix.exe -version
```

## Log:
To view the actual collected data, set the log to debug
```
.\processlist-OTE-0.79.1-exclude-thread-fix.exe -log=debug
```