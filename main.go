package main

import (
	_ "embed"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Usage = "aasaam analytics collector"
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		{
			Name:   "fixture",
			Usage:  "Run fixtures",
			Action: mainFixture,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "log-level",
					Usage:    "Could be one of `panic`, `fatal`, `error`, `warn`, `info`, `debug` or `trace`",
					Value:    "warn",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_LOG_LEVEL"},
				},

				&cli.StringFlag{
					Name:     "fixture-path",
					Usage:    "YAML fixtures path",
					Value:    "./fixture.yml",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_FIXTURE_PATH"},
				},
				&cli.IntFlag{
					Name:     "fixture-number",
					Usage:    "Number of each record try to insert",
					Value:    30,
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_FIXTURE_NUMBER"},
				},

				&cli.StringFlag{
					Name:     "redis-uri",
					Usage:    "Redis URI",
					Value:    "redis://127.0.0.1:6379/0",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_REDIS_URI"},
				},

				// clickhouse
				&cli.StringFlag{
					Name:     "clickhouse-servers",
					Usage:    "Comma separated clickhouse ip:port",
					Value:    "127.0.0.1:9440",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_CLICKHOUSE_SERVERS"},
				},
				&cli.StringFlag{
					Name:     "clickhouse-database",
					Usage:    "Clickhouse database name",
					Value:    "analytics",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_CLICKHOUSE_DATABASE"},
				},
				&cli.StringFlag{
					Name:     "clickhouse-username",
					Usage:    "Clickhouse username",
					Value:    "default",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_CLICKHOUSE_USERNAME"},
				},
				&cli.StringFlag{
					Name:     "clickhouse-password",
					Usage:    "Clickhouse password",
					Value:    "password123123",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_CLICKHOUSE_PASSWORD"},
				},
				&cli.IntFlag{
					Name:     "clickhouse-max-execution-time",
					Usage:    "Clickhouse max execution time",
					Value:    60,
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_CLICKHOUSE_MAX_EXECUTION_TIME"},
				},
				&cli.IntFlag{
					Name:     "clickhouse-dial-timeout",
					Usage:    "Clickhouse dial timeout in seconds",
					Value:    5,
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_CLICKHOUSE_DIAL_TIMEOUT"},
				},
				&cli.IntFlag{
					Name:     "clickhouse-max-idle-conns",
					Usage:    "Clickhouse max idle connections",
					Value:    5,
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_CLICKHOUSE_MAX_IDLE_CONNS"},
				},
				&cli.IntFlag{
					Name:     "clickhouse-max-open-conns",
					Usage:    "Clickhouse max open connections",
					Value:    10,
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_CLICKHOUSE_MAX_OPEN_CONNS"},
				},
				&cli.IntFlag{
					Name:     "clickhouse-conn-max-lifetime",
					Usage:    "Clickhouse connection max lifetime in seconds",
					Value:    3600,
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_CLICKHOUSE_CONN_MAX_LIFETIME"},
				},
				&cli.IntFlag{
					Name:     "clickhouse-max-block-size",
					Usage:    "Clickhouse max block size",
					Value:    10,
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_CLICKHOUSE_MAX_BLOCK_SIZE"},
				},
				&cli.BoolFlag{
					Name:     "clickhouse-compression-lz4",
					Usage:    "Clickhouse compression LZ4",
					Value:    false,
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_CLICKHOUSE_COMPRESSION_LZ4"},
				},
				&cli.StringFlag{
					Name:     "clickhouse-root-ca",
					Usage:    "Clickhouse Root CA path",
					Value:    "./cert/ca.pem",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_CLICKHOUSE_ROOT_CA"},
				},
				&cli.StringFlag{
					Name:     "clickhouse-client-cert",
					Usage:    "Clickhouse Client certificate path",
					Value:    "./cert/client-fullchain.pem",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_CLICKHOUSE_CLIENT_CERT"},
				},
				&cli.StringFlag{
					Name:     "clickhouse-client-key",
					Usage:    "Clickhouse Client key path",
					Value:    "./cert/client-key.pem",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_CLICKHOUSE_CLIENT_KEY"},
				},
			},
		},
		{
			Name:   "run-consumer",
			Usage:  "Run consumer to clickhouse",
			Action: mainStore,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "log-level",
					Usage:    "Could be one of `panic`, `fatal`, `error`, `warn`, `info`, `debug` or `trace`",
					Value:    "warn",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_LOG_LEVEL"},
				},

				&cli.StringFlag{
					Name:     "redis-uri",
					Usage:    "Redis URI",
					Value:    "redis://127.0.0.1:6379/0",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_REDIS_URI"},
				},
				&cli.IntFlag{
					Name:     "check-interval",
					Usage:    "Interval wait time",
					Value:    1,
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_CHECK_INTERVAL"},
				},

				// clickhouse
				&cli.StringFlag{
					Name:     "clickhouse-servers",
					Usage:    "Comma separated clickhouse ip:port",
					Value:    "127.0.0.1:9440",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_CLICKHOUSE_SERVERS"},
				},
				&cli.StringFlag{
					Name:     "clickhouse-database",
					Usage:    "Clickhouse database name",
					Value:    "analytics",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_CLICKHOUSE_DATABASE"},
				},
				&cli.StringFlag{
					Name:     "clickhouse-username",
					Usage:    "Clickhouse username",
					Value:    "default",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_CLICKHOUSE_USERNAME"},
				},
				&cli.StringFlag{
					Name:     "clickhouse-password",
					Usage:    "Clickhouse password",
					Value:    "password123123",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_CLICKHOUSE_PASSWORD"},
				},
				&cli.IntFlag{
					Name:     "clickhouse-max-execution-time",
					Usage:    "Clickhouse max execution time",
					Value:    60,
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_CLICKHOUSE_MAX_EXECUTION_TIME"},
				},
				&cli.IntFlag{
					Name:     "clickhouse-dial-timeout",
					Usage:    "Clickhouse dial timeout in seconds",
					Value:    5,
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_CLICKHOUSE_DIAL_TIMEOUT"},
				},
				&cli.IntFlag{
					Name:     "clickhouse-max-idle-conns",
					Usage:    "Clickhouse max idle connections",
					Value:    5,
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_CLICKHOUSE_MAX_IDLE_CONNS"},
				},
				&cli.IntFlag{
					Name:     "clickhouse-max-open-conns",
					Usage:    "Clickhouse max open connections",
					Value:    10,
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_CLICKHOUSE_MAX_OPEN_CONNS"},
				},
				&cli.IntFlag{
					Name:     "clickhouse-conn-max-lifetime",
					Usage:    "Clickhouse connection max lifetime in seconds",
					Value:    3600,
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_CLICKHOUSE_CONN_MAX_LIFETIME"},
				},
				&cli.IntFlag{
					Name:     "clickhouse-max-block-size",
					Usage:    "Clickhouse max block size",
					Value:    10,
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_CLICKHOUSE_MAX_BLOCK_SIZE"},
				},
				&cli.BoolFlag{
					Name:     "clickhouse-compression-lz4",
					Usage:    "Clickhouse compression LZ4",
					Value:    false,
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_CLICKHOUSE_COMPRESSION_LZ4"},
				},
				&cli.StringFlag{
					Name:     "clickhouse-root-ca",
					Usage:    "Clickhouse Root CA path",
					Value:    "./cert/ca.pem",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_CLICKHOUSE_ROOT_CA"},
				},
				&cli.StringFlag{
					Name:     "clickhouse-client-cert",
					Usage:    "Clickhouse Client certificate path",
					Value:    "./cert/client-fullchain.pem",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_CLICKHOUSE_CLIENT_CERT"},
				},
				&cli.StringFlag{
					Name:     "clickhouse-client-key",
					Usage:    "Clickhouse Client key path",
					Value:    "./cert/client-key.pem",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_CLICKHOUSE_CLIENT_KEY"},
				},
			},
		},
		{
			Name:   "run-producer",
			Usage:  "Run producer server listen on HTTP for receive data",
			Action: mainProducer,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "log-level",
					Usage:    "Could be one of `panic`, `fatal`, `error`, `warn`, `info`, `debug` or `trace`",
					Value:    "warn",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_LOG_LEVEL"},
				},

				&cli.StringFlag{
					Name:     "listen",
					Usage:    "Application listen http ip:port address",
					Value:    "0.0.0.0:4000",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_LISTEN"},
				},

				&cli.StringFlag{
					Name:     "collector-url",
					Usage:    "Full URL that expose collector url",
					Value:    "http://127.0.0.1:4000",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_COLLECTOR_URL"},
				},

				&cli.Int64Flag{
					Name:     "static-cache-ttl",
					Usage:    "Static cache max age for none versioning assets",
					Value:    8 * 3600,
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_STATIC_CACHE_TTL"},
				},

				&cli.StringFlag{
					Name:     "postgis-uri",
					Usage:    "Postgres geonames connection string",
					Value:    "postgres://geonames:geonames@127.0.0.1:5432/geonames",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_POSTGIS_URI"},
				},

				&cli.StringFlag{
					Name:     "redis-uri",
					Usage:    "Redis URI",
					Value:    "redis://127.0.0.1:6379/0",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_REDIS_URI"},
				},

				&cli.StringFlag{
					Name:     "mmdb-city-path",
					Usage:    "MMDB city database path",
					Value:    "tmp/GeoLite2-City.mmdb",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_MMDB_CITY_PATH"},
				},

				&cli.StringFlag{
					Name:     "mmdb-asn-path",
					Usage:    "MMDB asn database path",
					Value:    "tmp/GeoLite2-ASN.mmdb",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_MMDB_ASN_PATH"},
				},

				&cli.StringFlag{
					Name:     "management-projects-endpoint",
					Usage:    "URL of management server that expose projects",
					Value:    "http://localhost:9897/projects.json",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_MANAGEMENT_PROJECTS_ENDPOINT"},
				},

				&cli.StringFlag{
					Name:     "management-projects-json",
					Usage:    "Path of JSON file of projects",
					Value:    "",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_MANAGEMENT_PROJECTS_JSON"},
				},

				&cli.Int64Flag{
					Name:     "management-call-interval",
					Usage:    "Call update for projects in seconds",
					Value:    60,
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_MANAGEMENT_CALL_INTERVAL"},
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
