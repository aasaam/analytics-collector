package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/urfave/cli/v2"
)

func projectsLoad(url string) (map[string]projectData, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var r map[string]projectData
	errJSON := json.Unmarshal(body, &r)
	if errJSON != nil {
		return nil, err
	}
	return r, nil
}

func runServer(c *cli.Context) error {
	conf := newConfig(
		c.String("log-level"),
		c.Uint("static-cache-ttl"),
		c.Bool("script-source"),
		c.String("collector-url"),
		c.String("allowed-metrics-ips"),
	)

	st := newStorage()

	projectsManager := newProjectsManager()

	sleepTime := time.Duration(c.Int64("management-call-interval")) * time.Second

	go func() {

		for {
			prometheusUptimeInSeconds.Set(float64(time.Now().Unix() - initTime))

			projects, projectsErr := projectsLoad(c.String("management-projects-endpoint"))
			if projectsErr != nil {
				prometheusProjectsErrors.Inc()
				conf.getLogger().
					Error().
					Str("type", "projects_load").
					Str("on", "load_from_management").
					Str("error", projectsErr.Error()).
					Send()
				time.Sleep(sleepTime)
				continue
			}

			projectsManagerErr := projectsManager.load(projects)
			if projectsManagerErr != nil {
				prometheusProjectsErrors.Inc()
				conf.getLogger().
					Error().
					Str("type", "projects_load").
					Str("on", "load_data").
					Str("error", projectsManagerErr.Error()).
					Send()
				time.Sleep(sleepTime)
				continue
			}

			conf.getLogger().
				Trace().
				Str("type", "projects_load").
				Bool("success", true).
				Send()

			prometheusProjectsSuccess.Inc()

			time.Sleep(sleepTime)
		}
	}()

	ii := 0

	go func() {
		for {
			func() {

				records := st.getRecords()

				// fmt.Println("success")
				for _, i := range records {
					r := bytes.NewReader(i)

					var n2 record
					if err := gob.NewDecoder(r).Decode(&n2); err != nil {
						panic(n2)
					}
					fmt.Printf("%d %s\n", ii, n2.PURL.String())
					ii++
				}

				st.cleanRecords()

				// if st.count < 1 {
				// 	// fmt.Println("low storage")
				// 	time.Sleep(time.Duration(1) * time.Second)
				// 	return
				// }

				// t := time.Now().Unix() % 2

				//

				// if t == 1 { // success
				// 	// fmt.Println("success")
				// 	for _, i := range items {
				// 		// fmt.Fprintf(f, "%+v\n", i)
				// 		fmt.Printf("%s\n", i.PURL)
				// 	}

				// 	st.clean()
				// } else {
				// 	// fmt.Println("failed")
				// 	st.setItems(items)
				// }

			}()
			time.Sleep(time.Duration(1) * time.Second)
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

	app := newHTTPServer(conf, geoParser, projectsManager, st)
	return app.Listen(c.String("listen"))
}

func main() {
	app := cli.NewApp()
	app.Usage = "aasaam analytics collector"
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		{
			Name:   "run",
			Usage:  "Run server",
			Action: runServer,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "listen",
					Usage:    "Application listen http ip:port address",
					Value:    "0.0.0.0:4000",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_LISTEN_ADDRESS"},
				},
				&cli.StringFlag{
					Name:     "collector-url",
					Usage:    "Full URL that expose collector url",
					Value:    "http://127.0.0.1:4000",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_COLLECTOR_URL"},
				},
				&cli.StringFlag{
					Name:     "allowed-metrics-ips",
					Usage:    "Comma seprated ips that can access /metrics for prometheus exporter",
					Value:    "127.0.0.1",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_ALLOWED_METRICS_IPS"},
				},
				&cli.Int64Flag{
					Name:     "static-cache-ttl",
					Usage:    "Application listen http ip:port address",
					Value:    86400,
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_STATIC_CACHE_TTL"},
				},
				&cli.BoolFlag{
					Name:     "script-source",
					Usage:    "expose source of script for client debugging",
					Value:    false,
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_SCRIPT_SOURCE"},
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
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_MMDB_ASN_PATH"},
				},
				&cli.Int64Flag{
					Name:     "management-call-interval",
					Usage:    "Call update for projects in seconds",
					Value:    10,
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_MMDB_ASN_PATH"},
				},
				&cli.StringFlag{
					Name:     "log-level",
					Usage:    "Could be one of `panic`, `fatal`, `error`, `warn`, `info`, `debug` or `trace`",
					Value:    "warn",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_LOG_LEVEL"},
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
