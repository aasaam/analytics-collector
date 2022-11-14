# Installation

## 1- Copy Files
Copy ch0x directories to corresponding nodes


## 2- Generate password
```bash
  PASSWORD=$(base64 < /dev/urandom | head -c8); echo "$PASSWORD"; echo -n "$PASSWORD" | sha256sum | tr -d '-'
  ```
First line is your password and the second is `sha256_hex`, place the hash in the `password_sha256_hex` tag in the `config/users.xml` on each node:

```bash
<password_sha256_hex></password_sha256_hex>
```

## 3- Create .env:
On each node:
```bash
cp .env.sample .env
```

In the `.env` file change following Environment Variables according to your configuration:
```
CH_SERVER_HOST
CH_SERVER_OTHER2_HOST
CH_SERVER_OTHER3_HOST
CH_MAIN_DOMAIN
CH_CH01_IP
CH_CH02_IP
CH_CH03_IP
```

## 4- Create certificates:
Create certificates and copy generated files to `ch0x/config/certs` on each node.

## 5- Run clickhouse containers:
Run clickhouse containers on each server by following command:
```bash
docker-compose up -d
```

## 6- Run clickhouse client:
```bash
docker exec -ti ch01_server_1 clickhouse-client --user default --password [place plain password that generated in step 2] --port 9440 --secure --host 192.168.0.x
```

## 7- Import database schema
Import schema.sql file in each `clickhouse instance` line by line



**Your clickhouse cluster is ready to use**
