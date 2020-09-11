#!/bin/sh
cd webapp/go
go mod vendor
make build
cd ../..