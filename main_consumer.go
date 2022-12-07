package main

import (
	"time"

	"github.com/urfave/cli/v2"
)

func mainStore(c *cli.Context) error {
	conf := newConfig(
		c.String("log-level"),
		c.Uint("static-cache-ttl"),
		c.String("collector-url"),
	)

	clickhouseConfig := clickhouseConfig{
		servers:          c.String("clickhouse-servers"),
		database:         c.String("clickhouse-database"),
		username:         c.String("clickhouse-username"),
		password:         c.String("clickhouse-password"),
		maxExecutionTime: c.Int("clickhouse-max-execution-time"),
		dialTimeout:      c.Int("clickhouse-dial-timeout"),
		debug:            c.Bool("test-mode"),
		compressionLZ4:   c.Bool("clickhouse-compression-lz4"),
		maxIdleConns:     c.Int("clickhouse-max-idle-conns"),
		maxOpenConns:     c.Int("clickhouse-max-open-conns"),
		connMaxLifetime:  c.Int("clickhouse-conn-max-lifetime"),
		maxBlockSize:     c.Int("clickhouse-max-block-size"),
		rootCAPath:       c.String("clickhouse-root-ca"),
		clientCertPath:   c.String("clickhouse-client-cert"),
		clientKeyPath:    c.String("clickhouse-client-key"),
	}

	checkInterval := time.Duration(c.Int("check-interval")) * time.Second

	for {

		func() {

			r := workerRun(&clickhouseConfig, conf, c.String("redis-uri"))
			if r.e != nil {
				conf.getLogger().
					Error().
					Str("state", r.errorState).
					Str("error", r.e.Error()).
					Send()
			} else {
				conf.getLogger().
					Info().
					Int64("records", r.records).
					Int64("clientErrors", r.clientErrors).
					Float64("timeTaken", r.timeTaken).
					Send()
			}
		}()

		time.Sleep(checkInterval)
	}
}
