#!/bin/bash
cd /opt/isucon2020-himetex
git pull
sh tools/build.sh
sh tools/deploy-db.sh
sudo systemctl restart isuumo.go