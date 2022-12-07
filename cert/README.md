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

## Check Single

```bash
openssl s_client -CAfile ca.pem -cert client-fullchain.pem -key client-key.pem -connect ch.analytics-clickhouse.net.private:8443
openssl s_client -CAfile ca.pem -cert client-fullchain.pem -key client-key.pem -connect ch.analytics-clickhouse.net.private:9440
```

## Check Cluster

```bash
openssl s_client -CAfile ca.pem -cert client-fullchain.pem -key client-key.pem -connect ch1.analytics-clickhouse.net.private:8443
openssl s_client -CAfile ca.pem -cert client-fullchain.pem -key client-key.pem -connect ch1.analytics-clickhouse.net.private:9440
openssl s_client -CAfile ca.pem -cert client-fullchain.pem -key client-key.pem -connect ch1.analytics-clickhouse.net.private:9010
openssl s_client -CAfile ca.pem -cert client-fullchain.pem -key client-key.pem -connect ch1.analytics-clickhouse.net.private:9281
```
