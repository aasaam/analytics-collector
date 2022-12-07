#!/bin/bash

function validate_ip()
{
  local ip=$1
  local stat=1

  if [[ $ip =~ ^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$ ]]; then
    OIFS=$IFS
    IFS='.'
    ip=($ip)
    IFS=$OIFS
    [[ ${ip[0]} -le 255 && ${ip[1]} -le 255 \
    && ${ip[2]} -le 255 && ${ip[3]} -le 255 ]]
    stat=$?
  fi
  return $stat
}

CURRENT_DIR="$( dirname -- "$0"; )"
CURRENT_DIR="$( realpath -e -- "$CURRENT_DIR"; )";
PROJECT_DIR="$( realpath -e -- "$CURRENT_DIR/../"; )";

ASM_CH_NODE1_IP=$1
ASM_CH_NODE2_IP=$2
ASM_CH_NODE3_IP=$3
ASM_CH_MANGEGMENT_IP=$4
ASM_COLLECTOR_HOSTNAME=$5
ASM_MANEGMENT_HOSTNAME=$6

if [ -z "$ASM_COLLECTOR_HOSTNAME" ]; then
  echo "Collector hostname not set"
  exit 1
fi

if [ -z "$ASM_MANEGMENT_HOSTNAME" ]; then
  echo "Manegment hostname not set"
  exit 1
fi

if validate_ip $ASM_CH_NODE1_IP; then NODE01_IP_VALID='1'; else NODE01_IP_VALID='0'; fi
if validate_ip $ASM_CH_NODE2_IP; then NODE02_IP_VALID='1'; else NODE02_IP_VALID='0'; fi
if validate_ip $ASM_CH_NODE3_IP; then NODE03_IP_VALID='1'; else NODE03_IP_VALID='0'; fi
if validate_ip $ASM_CH_MANGEGMENT_IP; then MANGEGMENT_IP_VALID='1'; else MANGEGMENT_IP_VALID='0'; fi

if [[ $NODE01_IP_VALID == "0" ]]; then echo "invalid node 1 IP $ASM_CH_NODE1_IP"; exit 1; fi
if [[ $NODE02_IP_VALID == "0" ]]; then echo "invalid node 2 IP $ASM_CH_NODE2_IP"; exit 1; fi
if [[ $NODE03_IP_VALID == "0" ]]; then echo "invalid node 3 IP $ASM_CH_NODE3_IP"; exit 1; fi
if [[ $MANGEGMENT_IP_VALID == "0" ]]; then echo "invalid management IP $ASM_CH_MANGEGMENT_IP"; exit 1; fi

echo "All configurations seems fine:"
echo "Node 1 IP:            $ASM_CH_NODE1_IP"
echo "Node 2 IP:            $ASM_CH_NODE2_IP"
echo "Node 3 IP:            $ASM_CH_NODE3_IP"
echo "Magegment IP:         $ASM_CH_MANGEGMENT_IP"
echo "Collectors hostname:  $ASM_COLLECTOR_HOSTNAME"
echo "Magegment hostname:   $ASM_MANEGMENT_HOSTNAME"

CLICKHOUSE_PASSWORD=$(tr -dc A-Za-z0-9 </dev/urandom | head -c 48)
ASM_COLLECTOR_API_KEY=$(tr -dc A-Za-z0-9 </dev/urandom | head -c 48)
GRAFANA_PASSWORD=$(tr -dc A-Za-z0-9 </dev/urandom | head -c 48)
ASM_AUTH_HMAC_SECRET=$(openssl rand -hex 64)

MANEGMENT_PATH=$CURRENT_DIR/ready/management
mkdir -p $MANEGMENT_PATH
cp -rf $CURRENT_DIR/template/management/* $MANEGMENT_PATH/
cp -rf $CURRENT_DIR/template/management/.env $MANEGMENT_PATH/.env
cp -f $PROJECT_DIR/cert/{ca.pem,dhparam.pem,client-fullchain.pem,client-key.pem} $MANEGMENT_PATH/cert/
sed -i "s+__ASM_MANEGMENT_HOSTNAME__+$ASM_MANEGMENT_HOSTNAME+g" $MANEGMENT_PATH/.env
sed -i "s+__ASM_AUTH_HMAC_SECRET__+$ASM_AUTH_HMAC_SECRET+g" $MANEGMENT_PATH/.env
sed -i "s+__ASM_COLLECTOR_API_KEY__+$ASM_COLLECTOR_API_KEY+g" $MANEGMENT_PATH/.env
sed -i "s+__CLICKHOUSE_PASSWORD__+$CLICKHOUSE_PASSWORD+g" $MANEGMENT_PATH/.env
sed -i "s+__GRAFANA_PASSWORD__+$GRAFANA_PASSWORD+g" $MANEGMENT_PATH/.env
sed -i "s+__ASM_CH_NODE1_IP__+$ASM_CH_NODE1_IP+g" $MANEGMENT_PATH/.env
sed -i "s+__ASM_CH_NODE2_IP__+$ASM_CH_NODE2_IP+g" $MANEGMENT_PATH/.env
sed -i "s+__ASM_CH_NODE3_IP__+$ASM_CH_NODE3_IP+g" $MANEGMENT_PATH/.env
sed -i "s+__ASM_COLLECTOR_HOSTNAME__+$ASM_COLLECTOR_HOSTNAME+g" $MANEGMENT_PATH/.env

for i in $(seq 1 3); do
  NODE_PATH=$CURRENT_DIR/ready/node$i
  mkdir -p $NODE_PATH

  __NODE_ID__=$i
  __OTHER_NODE_1__="2"
  __OTHER_NODE_2__="3"
  if [[ $i == "2" ]]; then
    __OTHER_NODE_1__="1"
    __OTHER_NODE_2__="3"
  elif [[ $i == "3" ]]; then
    __OTHER_NODE_1__="2"
    __OTHER_NODE_2__="1"
  fi


  cp -rf $CURRENT_DIR/template/node/* $NODE_PATH/
  cp -f $PROJECT_DIR/cert/{ca.pem,dhparam.pem,client-fullchain.pem,client-key.pem,server-fullchain.pem,server-key.pem} $NODE_PATH/clickhouse/cert/
  cp -f $PROJECT_DIR/cert/{ca.pem,dhparam.pem,client-fullchain.pem,client-key.pem} $NODE_PATH/collector/cert/

  # clickhouse
  sed -i "s+__CLICKHOUSE_PASSWORD__+$CLICKHOUSE_PASSWORD+g" $NODE_PATH/clickhouse/.env
  sed -i "s+__ASM_CH_NODE1_IP__+$ASM_CH_NODE1_IP+g" $NODE_PATH/clickhouse/.env
  sed -i "s+__ASM_CH_NODE2_IP__+$ASM_CH_NODE2_IP+g" $NODE_PATH/clickhouse/.env
  sed -i "s+__ASM_CH_NODE3_IP__+$ASM_CH_NODE3_IP+g" $NODE_PATH/clickhouse/.env
  sed -i "s+__NODE_ID__+$__NODE_ID__+g" $NODE_PATH/clickhouse/.env
  sed -i "s+__OTHER_NODE_1__+$__OTHER_NODE_1__+g" $NODE_PATH/clickhouse/.env
  sed -i "s+__OTHER_NODE_2__+$__OTHER_NODE_2__+g" $NODE_PATH/clickhouse/.env

  rm -rf $NODE_PATH/clickhouse/nginx-exposer/acl.conf
  touch $NODE_PATH/clickhouse/nginx-exposer/acl.conf
  echo "allow $ASM_CH_NODE1_IP;" >> $NODE_PATH/clickhouse/nginx-exposer/acl.conf
  echo "allow $ASM_CH_NODE2_IP;" >> $NODE_PATH/clickhouse/nginx-exposer/acl.conf
  echo "allow $ASM_CH_NODE3_IP;" >> $NODE_PATH/clickhouse/nginx-exposer/acl.conf
  echo "allow $ASM_CH_MANGEGMENT_IP;" >> $NODE_PATH/clickhouse/nginx-exposer/acl.conf


  # collector
  sed -i "s+__NODE_ID__+$i+g" $NODE_PATH/collector/.env
  sed -i "s+__CLICKHOUSE_PASSWORD__+$CLICKHOUSE_PASSWORD+g" $NODE_PATH/collector/.env
  sed -i "s+__ASM_CH_NODE1_IP__+$ASM_CH_NODE1_IP+g" $NODE_PATH/collector/.env
  sed -i "s+__ASM_CH_NODE2_IP__+$ASM_CH_NODE2_IP+g" $NODE_PATH/collector/.env
  sed -i "s+__ASM_CH_NODE3_IP__+$ASM_CH_NODE3_IP+g" $NODE_PATH/collector/.env
  sed -i "s+__ASM_MANEGMENT_HOSTNAME__+$ASM_MANEGMENT_HOSTNAME+g" $NODE_PATH/collector/.env
  sed -i "s+__ASM_COLLECTOR_HOSTNAME__+$ASM_COLLECTOR_HOSTNAME+g" $NODE_PATH/collector/.env
  sed -i "s+__ASM_COLLECTOR_API_KEY__+$ASM_COLLECTOR_API_KEY+g" $NODE_PATH/collector/.env

  sed -i "s+__ASM_COLLECTOR_HOSTNAME__+$ASM_COLLECTOR_HOSTNAME+g" $NODE_PATH/collector/get-cloudflare-cert.sh
  chmod 500 $NODE_PATH/collector/get-cloudflare-cert.sh
  rm -rf $NODE_PATH/collector/tmp

done
