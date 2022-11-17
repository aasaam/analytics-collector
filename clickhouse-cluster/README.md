# Easy cluster setup

## Generate certificate

Using following steps:

1. Go to `clickhouse/cert`
2. Download latest [cfssl](https://github.com/cloudflare/cfssl):

   ```bash
   ./download-cfssl.sh
   ```

3. Generate certificates

   ```bash
   ./generate.sh
   ```

## Building cluster files

Using following steps:

1. Go to `clickhouse-cluster`
2. Initialize cluster according to your nodes IP:

   ```bash
   ./init.sh 192.168.56.201 192.168.56.202 192.168.56.203
   ```

3. Your nodes configuration files store on `clickhouse-cluster/ready`.
   Copy each node data to your desire path of server:

   ```bash
   scp -r ready/ch01 root@192.168.56.201:/root/
   scp -r ready/ch02 root@192.168.56.202:/root/
   scp -r ready/ch03 root@192.168.56.203:/root/
   ```

## Running clickhouse nodes

On each node run exact same commands:

```bash
cd /root/ch01
docker-compose up -d
docker exec -it analytics-clickhouse bash -c 'clickhouse-client --multiquery < /schema.sql'
```

You can browse using `clickhouse-client`:

```bash
docker exec -it analytics-clickhouse clickhouse-client --vertical -d analytics
```

**Note** Do not forget setup your ufw for firewall rules.
