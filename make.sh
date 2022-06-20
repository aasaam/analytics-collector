#!/bin/bash

set -e

# init
mkdir -p tmp
rm -rf embed/build
rm -rf tmp/analytics-client
mkdir -p embed/build

# mmdb
docker rm -f maxmind-lite-docker-test
docker run --name maxmind-lite-docker-test -d ghcr.io/aasaam/maxmind-lite-docker tail -f /dev/null
docker cp maxmind-lite-docker-test:/GeoLite2-City.mmdb tmp/GeoLite2-City.mmdb
docker cp maxmind-lite-docker-test:/GeoLite2-ASN.mmdb tmp/GeoLite2-ASN.mmdb
docker rm -f maxmind-lite-docker-test

# client
git clone --depth=1 https://github.com/aasaam/analytics-client.git tmp/analytics-client

cp tmp/analytics-client/dist/*.js embed/build/
cp tmp/analytics-client/dist/*.json embed/build/

wget -O embed/build/user_agents.yaml 'https://raw.githubusercontent.com/ua-parser/uap-core/master/regexes.yaml'
