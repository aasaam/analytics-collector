#!/bin/bash

set -e

if [ -f "server-fullchain.pem" ]; then
  echo "Certitifcates exists, try to remove them manually and run again."
  exit
fi

type ./cfssl >/dev/null 2>&1 || { echo >&2 "'cfssl' required. Aborting."; exit 1; }
type ./cfssljson >/dev/null 2>&1 || { echo >&2 "'cfssljson' required. Aborting."; exit 1; }

./cfssl gencert -initca csr-root.json | ./cfssljson -bare ca
./cfssl gencert -ca ca.pem -ca-key ca-key.pem -config ca-config.json -profile=server csr-server.json | ./cfssljson -bare server
./cfssl gencert -ca ca.pem -ca-key ca-key.pem -config ca-config.json -profile=client csr-client.json | ./cfssljson -bare client

openssl dhparam -out dhparam.pem 2048

cat server.pem ca.pem > server-fullchain.pem
cat client.pem ca.pem > client-fullchain.pem

chmod 644 *.pem
