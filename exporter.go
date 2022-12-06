package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var initTime = time.Now().Unix()

var (
	promMetricHTTPTotalRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "aasaam_analytics_collector_http_total_requests",
		Help: "Total number http requests",
	})

	promMetricHTTPErrors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "aasaam_analytics_collector_http_errors",
		Help: "Total number http response error",
	}, []string{"status"})

	promMetricInvalidRequestData = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "aasaam_analytics_collector_invalid_request_data",
		Help: "Number of invalid request data",
	}, []string{"on"})

	// promMetricValidRequestData = prometheus.NewCounterVec(prometheus.CounterOpts{
	// 	Name: "aasaam_analytics_collector_valid_request_data",
	// 	Help: "Number of valid request data",
	// }, []string{"type"})

	// promMetricInvalidProcessData = prometheus.NewGauge(prometheus.GaugeOpts{
	// 	Name: "aasaam_analytics_collector_invalid_process_data",
	// 	Help: "Number of invalid process data",
	// })

	promMetricRecordMode = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "aasaam_analytics_collector_record_mode",
		Help: "Total number record modes",
	}, []string{"mode"})

	// promMetricClientErrors = prometheus.NewCounter(prometheus.CounterOpts{
	// 	Name: "aasaam_analytics_collector_client_errors",
	// 	Help: "Total number of client errors",
	// })

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
	promRegistry.MustRegister(promMetricHTTPTotalRequests)
	promRegistry.MustRegister(promMetricHTTPErrors)
	promRegistry.MustRegister(promMetricRecordMode)
	// promRegistry.MustRegister(promMetricClientErrors)
	promRegistry.MustRegister(promMetricProjectsFetchErrors)
	promRegistry.MustRegister(promMetricProjectsFetchSuccess)
	promRegistry.MustRegister(promMetricInvalidRequestData)
	// promRegistry.MustRegister(promMetricInvalidProcessData)
	// promRegistry.MustRegister(promMetricValidRequestData)
	return promRegistry
}
