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

if validate_ip $ASM_CH_NODE1_IP; then NODE01_IP_VALID='1'; else NODE01_IP_VALID='0'; fi
if validate_ip $ASM_CH_NODE2_IP; then NODE02_IP_VALID='1'; else NODE02_IP_VALID='0'; fi
if validate_ip $ASM_CH_NODE3_IP; then NODE03_IP_VALID='1'; else NODE03_IP_VALID='0'; fi

if [[ $NODE01_IP_VALID == "0" ]]; then echo "invalid node 1 IP $ASM_CH_NODE1_IP"; exit 1; fi
if [[ $NODE02_IP_VALID == "0" ]]; then echo "invalid node 2 IP $ASM_CH_NODE2_IP"; exit 1; fi
if [[ $NODE03_IP_VALID == "0" ]]; then echo "invalid node 3 IP $ASM_CH_NODE3_IP"; exit 1; fi

echo "All nodes IP seems fine:"
echo "Node 1 IP: $ASM_CH_NODE1_IP"
echo "Node 2 IP: $ASM_CH_NODE2_IP"
echo "Node 3 IP: $ASM_CH_NODE3_IP"

RANDOM_PASSWORD=$(tr -dc A-Za-z0-9 </dev/urandom | head -c 32)

for i in $(seq 1 3); do
  declare "NODE_PATH"=$CURRENT_DIR/ready/ch0$i
  mkdir -p $NODE_PATH
  cp -rf $CURRENT_DIR/template/* $NODE_PATH/
  cp -f $PROJECT_DIR/clickhouse/cert/{ca.pem,dhparam.pem,client-fullchain.pem,client-key.pem,server-fullchain.pem,server-key.pem} $NODE_PATH/cert/
  cp -rf $CURRENT_DIR/template/.env $NODE_PATH/.env
  cp -rf $CURRENT_DIR/template/ufw.sh $NODE_PATH/ufw.sh
  sed -i "s+__RANDOM_PASSWORD__+$RANDOM_PASSWORD+g" $NODE_PATH/.env

  sed -i "s+__ASM_CH_NODE1_IP__+$ASM_CH_NODE1_IP+g" $NODE_PATH/.env
  sed -i "s+__ASM_CH_NODE2_IP__+$ASM_CH_NODE2_IP+g" $NODE_PATH/.env
  sed -i "s+__ASM_CH_NODE3_IP__+$ASM_CH_NODE3_IP+g" $NODE_PATH/.env

  sed -i "s+__ASM_CH_NODE1_IP__+$ASM_CH_NODE1_IP+g" $NODE_PATH/ufw.sh
  sed -i "s+__ASM_CH_NODE2_IP__+$ASM_CH_NODE2_IP+g" $NODE_PATH/ufw.sh
  sed -i "s+__ASM_CH_NODE3_IP__+$ASM_CH_NODE3_IP+g" $NODE_PATH/ufw.sh

  chmod +x $NODE_PATH/ufw.sh

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

  sed -i "s+__NODE_ID__+$__NODE_ID__+g" $NODE_PATH/.env
  sed -i "s+__OTHER_NODE_1__+$__OTHER_NODE_1__+g" $NODE_PATH/.env
  sed -i "s+__OTHER_NODE_2__+$__OTHER_NODE_2__+g" $NODE_PATH/.env

done

# NODE_1_PATH=$PROJECT_DIR/ready/ch01
# NODE_2_PATH=$PROJECT_DIR/ready/ch02
# NODE_3_PATH=$PROJECT_DIR/ready/ch03

# mkdir -p $NODE_1_PATH
# mkdir -p $NODE_2_PATH
# mkdir -p $NODE_3_PATH

# cp -f $PROJECT_DIR/../clickhouse/cert/{ca.pem,dhparam.pem,client-fullchain.pem,client-key.pem,server-fullchain.pem,server-key.pem} $NODE_1_PATH/cert/
# cp -f $PROJECT_DIR/../clickhouse/cert/{ca.pem,dhparam.pem,client-fullchain.pem,client-key.pem,server-fullchain.pem,server-key.pem} $NODE_2_PATH/cert/
# cp -f $PROJECT_DIR/../clickhouse/cert/{ca.pem,dhparam.pem,client-fullchain.pem,client-key.pem,server-fullchain.pem,server-key.pem} $NODE_3_PATH/cert/

# cp -rf $PROJECT_DIR/template/* $NODE_1_PATH/
# cp -rf $PROJECT_DIR/template/* $NODE_2_PATH/
# cp -rf $PROJECT_DIR/template/* $NODE_3_PATH/

# cat $PROJECT_DIR/template/.env |
