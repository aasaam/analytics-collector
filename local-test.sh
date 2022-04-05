#!/bin/bash

counter=0
while true
do
  # loop infinitely
  tmm=`date '+%Y%m%d%H%M%S.%N'`
  curl -s -H 'X-Forwarded-For: 46.4.106.18' "http://192.168.1.212:4000/?m=pv_ins&i=000000000000&u=https%3A%2F%2Fexample.net%2F?t=$tmm" -o /dev/null
  echo "$counter https://example.net/?t=$tmm"
  sleep 0.3
  ((counter++))
done
