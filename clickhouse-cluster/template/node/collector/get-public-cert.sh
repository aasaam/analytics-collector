#!/bin/bash

set -e

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

if [ -z "CLOUDFLARE_DNS_API_TOKEN" ]; then
  echo "CLOUDFLARE_DNS_API_TOKEN not defined"
  exit 1
fi

if [ -z "EMAIL" ]; then
  echo "EMAIL not defined"
  exit 1
fi

if [ ! -f ./lego ]; then
	echo "lego not found try to download..."
  mkdir -p $SCRIPT_DIR/tmp
  cd $SCRIPT_DIR/tmp
  curl -s https://api.github.com/repos/go-acme/lego/releases/latest \
  | jq -r '.assets[].browser_download_url' \
  | grep 'linux_amd64.tar.gz' \
  | wget -O lego.tar.gz -i -
  tar xf lego.tar.gz
  mv lego ../lego
fi

cd $SCRIPT_DIR

CLOUDFLARE_DNS_API_TOKEN="$CLOUDFLARE_DNS_API_TOKEN" ./lego --email "$EMAIL" --dns cloudflare --domains __ASM_COLLECTOR_DOMAIN__ --domains *.__ASM_COLLECTOR_DOMAIN__ run
cp ./.lego/certificates/__ASM_COLLECTOR_DOMAIN__.crt ./public-cert/fullchain.pem
cp ./.lego/certificates/__ASM_COLLECTOR_DOMAIN__.issuer.crt ./public-cert/chain.pem
cp ./.lego/certificates/__ASM_COLLECTOR_DOMAIN__.key ./public-cert/privkey.pem
openssl dhparam -out ./public-cert/dhparam.pem 2048
