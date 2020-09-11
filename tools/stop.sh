#!/bin/sh
ps aux | grep "isucon" | grep -v grep | awk '{ print "kill -9", $2 }' | sh