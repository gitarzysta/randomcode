#!/bin/sh
go build client.go
go build generate.go
go build server.go
mkdir -p data
echo Run ./generate to generate new data
