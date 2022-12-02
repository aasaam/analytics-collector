# Certificate

This is how to test your certificate:

## Single

```bash
openssl s_client -CAfile ca.pem -cert client-fullchain.pem -key client-key.pem -connect ch.analytics-clickhouse.net.private:8443
openssl s_client -CAfile ca.pem -cert client-fullchain.pem -key client-key.pem -connect ch.analytics-clickhouse.net.private:9440
```

## Cluster

```bash
openssl s_client -CAfile ca.pem -cert client-fullchain.pem -key client-key.pem -connect ch01.analytics-clickhouse.net.private:8443
openssl s_client -CAfile ca.pem -cert client-fullchain.pem -key client-key.pem -connect ch01.analytics-clickhouse.net.private:9440
openssl s_client -CAfile ca.pem -cert client-fullchain.pem -key client-key.pem -connect ch01.analytics-clickhouse.net.private:9010
openssl s_client -CAfile ca.pem -cert client-fullchain.pem -key client-key.pem -connect ch01.analytics-clickhouse.net.private:9281
```
