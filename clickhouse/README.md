# Running single mode clickhouse

Then create database and tables:

```bash
docker exec -it analytics-clickhouse bash -c 'clickhouse-client --multiquery < /schema.sql'
```

For using client:

```bash
docker exec -it analytics-clickhouse clickhouse-client --vertical -d analytics
```
