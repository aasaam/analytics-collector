#!/bin/bash
set -e

export CFSSL_VERSION='1.6.1'

curl -L "https://github.com/cloudflare/cfssl/releases/download/v${CFSSL_VERSION}/cfssl_${CFSSL_VERSION}_linux_amd64" -o cfssl
curl -L "https://github.com/cloudflare/cfssl/releases/download/v${CFSSL_VERSION}/cfssljson_${CFSSL_VERSION}_linux_amd64" -o cfssljson

chmod +x ./cfssl
chmod +x ./cfssljson
