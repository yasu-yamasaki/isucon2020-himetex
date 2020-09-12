#!/bin/bash
pushd /opt/isucon2020-himetex
git pull
sh tools/stop.sh
sh tools/build.sh
sudo systemctl start isuumo
popd