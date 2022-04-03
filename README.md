<div align="center">
  <h1>
    Analytics Collector
  </h1>
  <p>
    Analytics data collector and pre store processor for aasaam analytics
  </p>
  <p>
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
docker-compose -f docker-compose.dev.yml up -d
# import clickhouse schema
docker exec -t clickhouse-client /usr/bin/clickhouse-client --multiquery --host clickhouse-single --user 'analytics' --password 'password123123' < clickhouse/schema.sql

# run cli
docker exec -it clickhouse-client /usr/bin/clickhouse-client --host clickhouse-single --user 'analytics' --password 'password123123'
```

<div>
  <p align="center">
    <img alt="aasaam software development group" width="64" src="https://raw.githubusercontent.com/aasaam/information/master/logo/aasaam.svg">
    <br />
    aasaam software development group
  </p>
</div>
