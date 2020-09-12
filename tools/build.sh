#!/bin/sh
cd go
go mod vendor
make isuumo
rm /home/isucon/isuumo/webapp/go/isuumo
cp isuumo /home/isucon/isuumo/webapp/go/isuumo
cd ..