package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/urfave/cli/v2"
)

//go:embed clickhouse/insert_records.sql
var clickhouseInsertRecords string

//go:embed clickhouse/insert_client_errors.sql
var clickhouseInsertClientErrors string

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
		c.Bool("test-mode"),
		c.String("collector-url"),
		c.String("allowed-metrics-ips"),
	)

	clickhouseInit, _, clickhouseInitErr := getClickhouseConnection(
		c.String("clickhouse-servers"),
		c.String("clickhouse-database"),
		c.String("clickhouse-username"),
		c.String("clickhouse-password"),
		c.Int("clickhouse-max-execution-time"),
		c.Int("clickhouse-dial-timeout"),
		c.Bool("test-mode"),
		c.Bool("clickhouse-compression-lz4"),
		c.Int("clickhouse-max-idle-conns"),
		c.Int("clickhouse-max-open-conns"),
		c.Int("clickhouse-conn-max-lifetime"),
		c.Int("clickhouse-max-block-size"),
		nil,
		nil,
	)

	/*
		progress func(p *clickhouse.Progress),
		profile func(p *clickhouse.ProfileInfo),
	*/
	if clickhouseInitErr != nil {
		return clickhouseInitErr
	}

	clickhouseInit.Close()
	conf.getLogger().
		Debug().
		Msg("successfully ping to clickhouse")

	storage := newStorage()
	projectsManager := newProjectsManager()

	/**
	 * Projects
	 */
	sleepTime := time.Duration(c.Int64("management-call-interval")) * time.Second
	go func() {
		for {
			promMetricUptimeInSeconds.Set(float64(time.Now().Unix() - initTime))

			projects, projectsErr := projectsLoad(c.String("management-projects-endpoint"))
			if projectsErr != nil {
				promMetricProjectsFetchErrors.Inc()
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
				promMetricProjectsFetchErrors.Inc()
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

			promMetricProjectsFetchSuccess.Inc()
			time.Sleep(sleepTime)
		}
	}()

	/**
	 * Records
	 */
	clickhouseInterval := time.Duration(c.Int("clickhouse-interval"))
	go func() {
		for {
			func() {
				storage.Lock()
				defer storage.Unlock()

				// no storage data check
				if storage.recordCount == 0 && storage.clientErrorCount == 0 {
					conf.getLogger().
						Debug().
						Msg("storage is empty")
					time.Sleep(clickhouseInterval * time.Second)
					return
				}

				clickhouseConn, clickhouseCtx, clickhouseConnErr := getClickhouseConnection(
					c.String("clickhouse-servers"),
					c.String("clickhouse-database"),
					c.String("clickhouse-username"),
					c.String("clickhouse-password"),
					c.Int("clickhouse-max-execution-time"),
					c.Int("clickhouse-dial-timeout"),
					c.Bool("test-mode"),
					c.Bool("clickhouse-compression-lz4"),
					c.Int("clickhouse-max-idle-conns"),
					c.Int("clickhouse-max-open-conns"),
					c.Int("clickhouse-conn-max-lifetime"),
					c.Int("clickhouse-max-block-size"),
					nil,
					nil,
				)

				if clickhouseConnErr != nil {
					conf.getLogger().
						Error().
						Str("type", errorTypeApp).
						Str("on", "clickhouse-connection").
						Str("error", clickhouseConnErr.Error()).
						Send()
					time.Sleep(clickhouseInterval * time.Second)
					return
				}

				//
				// records
				//
				if storage.recordCount > 0 {
					records := storage.getRecords()

					recordsBatch, recordsBatchErr := clickhouseConn.PrepareBatch(
						clickhouseCtx, clickhouseInsertRecords,
					)
					if recordsBatchErr != nil {
						conf.getLogger().
							Error().
							Str("type", errorTypeApp).
							Str("on", "clickhouse-connection").
							Str("error", recordsBatchErr.Error()).
							Send()
						time.Sleep(clickhouseInterval * time.Second)
						return
					}

					for _, recordByte := range records {
						recordByteReader := bytes.NewReader(recordByte)

						var rec record
						recordDecodeErr := gob.NewDecoder(recordByteReader).Decode(&rec)
						if recordDecodeErr != nil {
							conf.getLogger().
								Error().
								Str("type", errorTypeApp).
								Str("on", "record-decode").
								Str("error", recordDecodeErr.Error()).
								Send()
							continue
						}

						if rec.EventCount > 0 {
							for i := 0; i < rec.EventCount; i++ {
								ECategory := rec.Events[i].ECategory
								EAction := rec.Events[i].EAction
								ELabel := rec.Events[i].ELabel
								EValue := rec.Events[i].EValue
								insertErr := insertRecordBatch(recordsBatch, rec, ECategory, EAction, ELabel, EValue)
								if insertErr != nil {
									conf.getLogger().
										Error().
										Str("type", errorTypeApp).
										Str("on", "record-insert").
										Str("error", insertErr.Error()).
										Send()
								}
							}
						} else {
							insertErr := insertRecordBatch(recordsBatch, rec, "", "", "", 0)
							if insertErr != nil {
								conf.getLogger().
									Error().
									Str("type", errorTypeApp).
									Str("on", "record-insert").
									Str("error", insertErr.Error()).
									Send()
							}
						}
					}

					recordsBatchSendErr := recordsBatch.Send()
					if recordsBatchSendErr != nil {
						conf.getLogger().
							Error().
							Str("type", errorTypeApp).
							Str("on", "record-batch-send").
							Str("error", recordsBatchSendErr.Error()).
							Send()
					}

					storage.cleanRecords()
				}

				//
				// client errors
				//
				if storage.clientErrorCount > 0 {
					clientErrors := storage.getClientErrors()

					clientErrorsBatch, clientErrorsBatchErr := clickhouseConn.PrepareBatch(
						clickhouseCtx, clickhouseInsertClientErrors,
					)
					if clientErrorsBatchErr != nil {
						conf.getLogger().
							Error().
							Str("type", errorTypeApp).
							Str("on", "clickhouse-connection").
							Str("error", clientErrorsBatchErr.Error()).
							Send()
						time.Sleep(clickhouseInterval * time.Second)
						return
					}

					for _, clientErrorByte := range clientErrors {
						clientErrorByteReader := bytes.NewReader(clientErrorByte)

						var ce record
						clientErrorDecodeErr := gob.NewDecoder(clientErrorByteReader).Decode(&ce)
						if clientErrorDecodeErr != nil {
							conf.getLogger().
								Error().
								Str("type", errorTypeApp).
								Str("on", "client-error-decode").
								Str("error", clientErrorDecodeErr.Error()).
								Send()
							continue
						}

						insertErr := insertClientErrBatch(clientErrorsBatch, ce)
						if insertErr != nil {
							conf.getLogger().
								Error().
								Str("type", errorTypeApp).
								Str("on", "client-error-insert").
								Str("error", insertErr.Error()).
								Send()
						}
					}

					clientErrorsBatchSendErr := clientErrorsBatch.Send()
					if clientErrorsBatchSendErr != nil {
						conf.getLogger().
							Error().
							Str("type", errorTypeApp).
							Str("on", "client-error-batch-send").
							Str("error", clientErrorsBatchSendErr.Error()).
							Send()
					}

					storage.cleanRecords()
				}
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

	app := newHTTPServer(conf, geoParser, projectsManager, storage)
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
					Value:    "debug",
					Required: false,
					EnvVars:  []string{"ASM_ANALYTICS_COLLECTOR_LOG_LEVEL"},
				},
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
