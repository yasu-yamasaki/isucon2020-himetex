#!/bin/sh
cd go
go mod vendor
make isuumo
cp isuumo /home/isucon/isuumo/webapp/go/isuumo
cd ..