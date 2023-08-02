# Build:
```
GOARCH=amd64 GOOS=windows GO111MODULE=on CGO_ENABLED=0 go build -o ./processlist-SApkgs-5.27.3.exe
```

# Usage:
## Help:
```
.\processlist-otelpkgs-0.79.1.exe -h
Usage of C:\Users\Administrator\Downloads\processlist-SApkgs-5.27.3.exe:
  -log string
        debug : Display additional output (default "info")
  -version
        Display version

```
## Version:
To view the module versions, run it with --version
```
.\processlist-SApkgs-5.27.3.exe -version
```

## Log:
To view the actual collected data, set the log to debug
```
.\processlist-SApkgs-5.27.3.exe -log=debug
```
