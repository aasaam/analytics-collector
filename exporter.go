package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var initTime = time.Now().Unix()

var (
	prometheusUptimeInSeconds = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "analytics_collector_uptime_in_seconds",
		Help: "Uptime seconds during application running",
	})

	prometheusStorageRecords = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "analytics_collector_storage_records",
		Help: "Number of storage records",
	})

	prometheusStorageClientErrors = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "analytics_collector_storage_client_errors",
		Help: "Number of storage client errors",
	})

	prometheusTotalRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "analytics_collector_total_requests",
		Help: "The total number requests",
	})

	prometheusClientErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "analytics_collector_client_errors",
		Help: "The total number of client errors",
	})

	prometheusProjectsErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "analytics_collector_projects_fetch_errors",
		Help: "The total number of projects fetch errors",
	})

	prometheusProjectsSuccess = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "analytics_collector_projects_fetch_success",
		Help: "The total number of projects fetch success",
	})

	prometheusResponseErrors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "analytics_collector_response_errors",
		Help: "The total number of response error",
	}, []string{"status"})
)

func getPrometheusRegistry() *prometheus.Registry {
	promRegistry := prometheus.NewRegistry()
	promRegistry.MustRegister(prometheusClientErrors)
	promRegistry.MustRegister(prometheusProjectsErrors)
	promRegistry.MustRegister(prometheusProjectsSuccess)
	promRegistry.MustRegister(prometheusResponseErrors)
	promRegistry.MustRegister(prometheusStorageClientErrors)
	promRegistry.MustRegister(prometheusStorageRecords)
	promRegistry.MustRegister(prometheusTotalRequests)
	promRegistry.MustRegister(prometheusUptimeInSeconds)

	prometheusUptimeInSeconds.Set(0)
	return promRegistry
}
