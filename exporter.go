package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var initTime = time.Now().Unix()

var (
	promMetricUptimeInSeconds = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "aasaam_analytics_collector_uptime_in_seconds",
		Help: "Uptime in seconds that collector is running",
	})

	promMetricStorageQueueRecords = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "aasaam_analytics_collector_storage_queue_records",
		Help: "Storage number of records that not inserted yet",
	})

	promMetricStorageQueueClientErrors = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "aasaam_analytics_collector_storage_queue_client_errors",
		Help: "Storage number of client errors that not inserted yet",
	})

	promMetricHTTPTotalRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "aasaam_analytics_collector_http_total_requests",
		Help: "Total number http requests",
	})

	promMetricHTTPErrors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "aasaam_analytics_collector_http_errors",
		Help: "Total number http response error",
	}, []string{"status"})

	promMetricRecordMode = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "aasaam_analytics_collector_record_mode",
		Help: "Total number record modes",
	}, []string{"mode"})

	promMetricClientErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "aasaam_analytics_collector_client_errors",
		Help: "Total number of client errors",
	})

	promMetricProjectsFetchErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "aasaam_analytics_collector_projects_fetch_errors",
		Help: "Total number of projects fetch errors",
	})

	promMetricProjectsFetchSuccess = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "aasaam_analytics_collector_projects_fetch_success",
		Help: "Total number of projects fetch success",
	})
)

func getPrometheusRegistry() *prometheus.Registry {
	promRegistry := prometheus.NewRegistry()
	promRegistry.MustRegister(promMetricUptimeInSeconds)
	promRegistry.MustRegister(promMetricStorageQueueRecords)
	promRegistry.MustRegister(promMetricStorageQueueClientErrors)
	promRegistry.MustRegister(promMetricHTTPTotalRequests)
	promRegistry.MustRegister(promMetricHTTPErrors)
	promRegistry.MustRegister(promMetricRecordMode)
	promRegistry.MustRegister(promMetricClientErrors)
	promRegistry.MustRegister(promMetricProjectsFetchErrors)
	promRegistry.MustRegister(promMetricProjectsFetchSuccess)
	return promRegistry
}
