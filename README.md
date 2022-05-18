# Tiny File Server

## What's this?

This tiny file server works in a directory that you specified and provide file service including file downloading and file uploading. 
By exposing its ip and listening port on the Internet, you can access your file anytime and anywhere with no speed limits.

## Quick Start

Create a configuration yaml file under your home path named with `.fileserver.yaml` to specify the server's listening port and base working directory.

```yaml
addr : ":11111"
basePath : "/home/fileserver/shared-data/"
```

Run the commands below to run set up your own server.

```shell
go build -o fileserver main.go
./fileserver
```

Also, you can make it a service and start it on booting.

## Attention

So far this server has no authentication and once you deploy it on the Internet, anyone has access to your file with full permission. This is very dangerous.

Its UI is really poor, and you're free to modify it. Have fun :)