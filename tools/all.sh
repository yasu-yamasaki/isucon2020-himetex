#!/bin/bash
cd /opt/isucon2020-himetex
git pull
sh tools/build.sh
sudo systemctl restart isuumo.go