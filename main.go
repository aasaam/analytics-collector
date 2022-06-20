package main

import (
	_ "embed"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

// //go:embed clickhouse/insert_records.sql
// var clickhouseInsertRecords string

// //go:embed clickhouse/insert_client_errors.sql
// var clickhouseInsertClientErrors string

func main() {
	app := cli.NewApp()
	app.Usage = "aasaam analytics collector"
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		{
			Name:   "run",
			Usage:  "Run collect server",
			Action: mainRun,
			Flags: []cli.Flag{
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
					Usage:    "Application listen http ip:port address",
					Value:    4 * 3600,
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_STATIC_CACHE_TTL"},
				},
				&cli.BoolFlag{
					Name:     "test-mode",
					Usage:    "Enable test mode",
					Value:    false,
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_TEST_MODE"},
				},
				&cli.StringFlag{
					Name:     "postgis-uri",
					Usage:    "Postgres geonames connection string",
					Value:    "postgres://geonames:geonames@127.0.0.1:5432/geonames",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_POSTGIS_URI"},
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
					Name:     "redis-queue-size",
					Usage:    "Max redis queue size for each insert",
					Value:    1000,
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_REDIS_QUEUE_SIZE"},
				},
				&cli.Int64Flag{
					Name:     "management-call-interval",
					Usage:    "Call update for projects in seconds",
					Value:    10,
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_MANAGEMENT_CALL_INTERVAL"},
				},
				&cli.StringFlag{
					Name:     "log-level",
					Usage:    "Could be one of `panic`, `fatal`, `error`, `warn`, `info`, `debug` or `trace`",
					Value:    "debug",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_LOG_LEVEL"},
				},
				// clickhouse
				&cli.IntFlag{
					Name:     "clickhouse-interval",
					Usage:    "Clickhouse interval in seconds",
					Value:    5,
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_CLICKHOUSE_INTERVAL"},
				},
				&cli.StringFlag{
					Name:     "clickhouse-servers",
					Usage:    "Comma separeted clickhouse ip:port",
					Value:    "127.0.0.1:9000",
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
					Value:    "analytics",
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
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
