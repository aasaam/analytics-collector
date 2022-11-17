#!/bin/bash

declare -a SERVER_IPS=("__ASM_CH_NODE1_IP__" "__ASM_CH_NODE2_IP__" "__ASM_CH_NODE3_IP__")
declare -a PUBLIC_PORTS=("8443" "9440")
declare -a INTERNAL_PORTS=("9010" "9281" "9234")

for SERVER_IP in "${SERVER_IPS[@]}"
do
  for INTERNAL_PORT in "${INTERNAL_PORTS[@]}"
  do
    ufw allow from $SERVER_IP to any port $INTERNAL_PORT comment "clickhouse private port"
  done
done

for PUBLIC_PORT in "${PUBLIC_PORTS[@]}"
do
  ufw allow $PUBLIC_PORT/tcp comment "clickhouse public port"
done
