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
			Name:   "run-collector",
			Usage:  "Run collect server listen on HTTP for receive data",
			Action: mainRun,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:     "test-mode",
					Usage:    "Enable test mode",
					Value:    false,
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_TEST_MODE"},
				},
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
