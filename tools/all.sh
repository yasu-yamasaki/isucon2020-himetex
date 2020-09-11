#!/bin/bash
git pull
sh tools/stop.sh
sh tools/build.sh
sudo systemctl start iscon