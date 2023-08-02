# Build:
```
GOARCH=amd64 GOOS=windows GO111MODULE=on CGO_ENABLED=0 go build -o ./processlist-otelpkgs-0.79.1.exe
```

# Usage:
## Help:
```
.\processlist-otelpkgs-0.79.1.exe -h
Usage of C:\Users\Administrator\Downloads\processlist-otelpkgs-0.79.1.exe:
  -log string
        debug : Display additional output (default "info")
  -version
        Display version

```
## Version:
To view the module versions, run it with --version
```
.\processlist-otelpkgs-0.79.1.exe -version
```

## Log:
To view the actual collected data, set the log to debug
```
.\processlist-otelpkgs-0.79.1.exe -log=debug
```
