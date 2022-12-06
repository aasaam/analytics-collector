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
		c.String("collector-url"),
	)

	managementCallInterval := time.Duration(c.Int64("management-call-interval")) * time.Second

redisStep:

	red, redErr := redisGetClient(c.String("redis-uri"))

	if redErr != nil {
		conf.getLogger().
			Error().
			Msg(redErr.Error())
		time.Sleep(time.Duration(1) * time.Second)
		goto redisStep
	}

	conf.getLogger().
		Debug().
		Msg("successfully ping to redis")

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
		red,
	)

	return app.Listen(c.String("listen"))
}
