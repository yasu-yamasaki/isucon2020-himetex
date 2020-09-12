#!/bin/sh
ps aux | grep "isuumo" | grep -v grep | awk '{ print "kill -9", $2 }' | sh