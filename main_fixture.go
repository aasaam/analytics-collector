package main

import (
	"context"
	"time"

	"github.com/urfave/cli/v2"
)

func mainFixture(c *cli.Context) error {
	conf := newConfig(
		c.String("log-level"),
		c.Uint("static-cache-ttl"),
		c.String("collector-url"),
	)

	numberOfFixtures := intMinMax(c.Int("fixture-number"), 1, 500)

	conf.getLogger().
		Info().
		Int("fixture-number", numberOfFixtures).
		Msg("Number of records on each cycle")

	f, fE := fixtureLoad(c.String("fixture-path"))
	if fE != nil {
		return fE
	}

redisStep:

	redisClient, redisClientErr := redisGetClient(c.String("redis-uri"))

	if redisClientErr != nil {
		conf.getLogger().
			Error().
			Msg(redisClientErr.Error())
		time.Sleep(time.Duration(1) * time.Second)
		goto redisStep
	}

	conf.getLogger().
		Debug().
		Msg("successfully ping to redis")

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

	refererParser := newRefererParser()
	userAgentParser := newUserAgentParser()

	for {
		go func() {

			for i := 0; i <= numberOfFixtures; i++ {
				r := f.record(refererParser, userAgentParser)
				rb, rbE := r.finalize()
				if rbE == nil {
					_, redErr := redisClient.LPush(context.Background(), redisKeyRecords, rb).Result()
					if redErr != nil {
						conf.getLogger().
							Error().
							Str("error", redErr.Error()).
							Send()
						return
					}
				}

			}

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

		time.Sleep(1 * time.Second)
	}
}
