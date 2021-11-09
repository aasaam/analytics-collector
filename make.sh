#!/bin/bash

# init
rm -rf tmp
mkdir tmp
rm -rf embed/build
mkdir embed/build

# client
git clone --depth=1 https://github.com/aasaam/analytics-client.git tmp/analytics-client

cp tmp/analytics-client/dist/amp.json embed/build/amp.json
cp tmp/analytics-client/dist/index.js embed/build/index.js
cp tmp/analytics-client/dist/iframe.js embed/build/iframe.js

wget -O embed/build/user_agents.yaml 'https://raw.githubusercontent.com/ua-parser/uap-core/master/regexes.yaml'
