# Build:
```
GOARCH=amd64 GOOS=windows GO111MODULE=on CGO_ENABLED=0 go build -o ./processlist-SA-5.0.2.exe
```

# Usage:
## Help:
```
.\processlist-SA-5.0.2.exe -h
Usage of C:\Users\Administrator\Downloads\processlist-SA-5.0.2.exe:
  -log string
        debug : Display additional output (default "info")
  -version
        Display version

```
## Version:
To view the module versions, run it with --version
```
.\processlist-SA-5.0.2.exe -version
```

## Log:
To view the actual collected data, set the log to debug
```
.\processlist-SA-5.0.2.exe -log=debug
```
