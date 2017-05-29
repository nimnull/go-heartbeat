#!/bin/bash

netstat -npt | grep ':3000' | grep ESTABLISHED | awk '{ print $5 }' | awk -F ':' '{ print $1 }' | sort -u | wc -l
