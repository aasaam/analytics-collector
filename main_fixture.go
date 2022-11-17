package main

import (
	"time"

	"github.com/urfave/cli/v2"
)

func mainFixture(c *cli.Context) error {
	conf := newConfig(
		c.String("log-level"),
		c.Uint("static-cache-ttl"),
		c.Bool("test-mode"),
		c.String("collector-url"),
	)

	numberOfFixtures := c.Int("fixture-number")
	if numberOfFixtures < 1 {
		numberOfFixtures = 1
	} else if numberOfFixtures > 100 {
		numberOfFixtures = 100
	}

	conf.getLogger().
		Info().
		Int("fixture-number", numberOfFixtures).
		Msg("Number of records on each cycle")

	f, fE := fixtureLoad(c.String("fixture-path"))
	if fE != nil {
		return fE
	}

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

	clickhouseInterval := time.Duration(c.Int("clickhouse-interval")) * time.Second

clickHouseInitStep:

	clickhouseInit, _, clickhouseInitErr := clickhouseGetConnection(&clickhouseConfig)

	if clickhouseInitErr != nil {
		conf.getLogger().
			Error().
			Msg(clickhouseInitErr.Error())
		time.Sleep(clickhouseInterval)
		goto clickHouseInitStep
	}

	clickhouseInit.Close()
	conf.getLogger().
		Debug().
		Msg("successfully ping to clickhouse")

	/**
	 * storage
	 */
	storage := newStorage()

	refererParser := newRefererParser()
	userAgentParser := newUserAgentParser()

	for {
		go func() {

			for i := 0; i <= numberOfFixtures; i++ {
				r := f.record(refererParser, userAgentParser)
				rb, rbE := r.finalize()
				if rbE == nil {
					storage.addRecord(rb)
				}
			}

			r := workerRun(&clickhouseConfig, conf, storage)
			if r.e != nil {
				conf.getLogger().
					Error().
					Str("type", errorTypeApp).
					Str("state", r.errorState).
					Str("error", r.e.Error()).
					Send()
			} else {
				conf.getLogger().
					Info().
					Int64("records", r.records).
					Int64("clientErrors", r.clientErrors).
					Int64("timeTaken", r.timeTaken).
					Send()
			}

		}()
		time.Sleep(1 * time.Second)
	}
}
