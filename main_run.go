package main

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/urfave/cli/v2"
)

func mainRun(c *cli.Context) error {
	conf := newConfig(
		c.String("log-level"),
		c.Uint("static-cache-ttl"),
		c.Bool("test-mode"),
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

	managementCallInterval := time.Duration(c.Int64("management-call-interval")) * time.Second

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

	go func() {
		for {
			promMetricUptimeInSeconds.Set(float64(time.Now().Unix() - initTime))
			time.Sleep(1 * time.Second)
		}
	}()

	clickhouseInit.Close()
	conf.getLogger().
		Debug().
		Msg("successfully ping to clickhouse")

	/**
	 * storage
	 */
	storage := newStorage()

	/**
	 * projects
	 */
	projectsManager := newProjectsManager()
	projectsJSON := c.String("management-projects-json")

	// it's static from json file
	if projectsJSON != "" {

		projects, projectsErr := projectsLoadJSON(c.String("management-projects-json"))
		if projectsErr != nil {
			return projectsErr
		}

		projectsManagerErr := projectsManager.load(projects)
		if projectsManagerErr != nil {
			return projectsManagerErr
		}

	} else { // from management

		projects, projectsErr := projectsLoad(c.String("management-projects-endpoint"))
		if projectsErr != nil {
			return projectsErr
		}

		projectsManagerErr := projectsManager.load(projects)
		if projectsManagerErr != nil {
			return projectsManagerErr
		}

		go func() {
			for {
				time.Sleep(managementCallInterval)

				e := workerProjects(c.String("management-projects-endpoint"), projectsManager)
				if e != nil {
					promMetricProjectsFetchErrors.Inc()
					conf.getLogger().
						Error().
						Str("type", "projects_load").
						Str("e", e.Error()).
						Bool("success", false).
						Send()
					continue
				}

				conf.getLogger().
					Info().
					Str("type", "projects_load").
					Int("number", projectsManager.total).
					Bool("success", true).
					Send()

				promMetricProjectsFetchSuccess.Inc()
			}
		}()
	}

	go func() {
		for {

			time.Sleep(clickhouseInterval)
			func() {
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
						Int64("timeTakenMS", r.timeTakenMS).
						Send()
				}
			}()
		}
	}()

	conn, connErr := pgx.Connect(context.Background(), c.String("postgis-uri"))
	if connErr != nil {
		return connErr
	}

	defer conn.Close(context.Background())

	geoParser, geoParserErr := newGeoParser(conn, c.String("mmdb-city-path"), c.String("mmdb-asn-path"))
	if geoParserErr != nil {
		return geoParserErr
	}

	refererParser := newRefererParser()
	userAgentParser := newUserAgentParser()

	app := newHTTPServer(
		conf,
		geoParser,
		refererParser,
		userAgentParser,
		projectsManager,
		storage,
	)

	return app.Listen(c.String("listen"))
}
