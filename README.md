<div align="center">
  <h1>
    Analytics Collector
  </h1>
  <p>
    Analytics data collector and pre store processor for aasaam analytics
  </p>
  <p>
    <a href="https://github.com/aasaam/analytics-collector/actions/workflows/test.yml" target="_blank">
      <img src="https://github.com/aasaam/analytics-collector/actions/workflows/test.yml/badge.svg" alt="test" />
    </a>
    <a href="https://goreportcard.com/report/github.com/aasaam/analytics-collector">
      <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/aasaam/analytics-collector">
      </a>
    <a href="https://github.com/aasaam/analytics-collector/blob/master/LICENSE">
      <img alt="License" src="https://img.shields.io/github/license/aasaam/analytics-collector">
    </a>
  </p>
</div>

## Development

```bash
# prepare dependencies
./make.sh
sudo rm -rf clickhouse/db clickhouse/logs
docker-compose -f docker-compose.dev.yml up -d
# import clickhouse schema
docker exec -it clickhouse-client bash
clickhouse-client --multiquery --host clickhouse-single --user 'analytics' --password 'password123123' < /tmp/schema.sql

# run cli
docker exec -it clickhouse-client /usr/bin/clickhouse-client --vertical --host clickhouse-single --user 'analytics' --password 'password123123'
```

<div>
  <p align="center">
    <img alt="aasaam software development group" width="64" src="https://raw.githubusercontent.com/aasaam/information/master/logo/aasaam.svg">
    <br />
    aasaam software development group
  </p>
</div>
