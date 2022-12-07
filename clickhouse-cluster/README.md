# Easy cluster setup

## Certificate

This is how to test generate and test certificate

```bash
cd cert
```

### Download CFSSL

```bash
./download-cfssl.sh
```

### Generate

```bash
./generate.sh
```

## Building cluster files

Using following steps:

1. Go to `clickhouse-cluster`
2. Initialize cluster according to your nodes IP:

   ```bash
   ./create.sh 192.168.56.201 192.168.56.202 192.168.56.203 192.168.56.100 collector.tld management.your-company.tld
   ```

3. Your nodes configuration files store on `clickhouse-cluster/ready`.
   Copy each node data to your desire path of server:

   ```bash
   scp -r ready/node1 root@192.168.56.201:/root/
   scp -r ready/node2 root@192.168.56.202:/root/
   scp -r ready/node3 root@192.168.56.203:/root/
   scp -r ready/management root@192.168.56.100:/root/
   ```

## Running clickhouse nodes

On each node run exact same commands:

```bash
cd /root/node1/clickhouse
docker-compose up -d
docker exec -it analytics-clickhouse bash -c 'clickhouse-client --multiquery < /schema.sql'
```

You can browse using `clickhouse-client`:

```bash
docker exec -it analytics-clickhouse clickhouse-client --vertical --database analytics
```

## Running collector

On each node run exact same commands:

```bash
cd /root/node1/collector
./get-cloudflare-cert.sh
```

```bash
docker-compose up -d
```
