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

	promMetricValidRequestData = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "aasaam_analytics_collector_valid_request_data",
		Help: "Number of valid request data",
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
	promRegistry.MustRegister(promMetricHTTPTotalRequests)
	promRegistry.MustRegister(promMetricHTTPErrors)
	promRegistry.MustRegister(promMetricInvalidRequestData)
	promRegistry.MustRegister(promMetricValidRequestData)
	promRegistry.MustRegister(promMetricClientErrors)
	promRegistry.MustRegister(promMetricProjectsFetchErrors)
	promRegistry.MustRegister(promMetricProjectsFetchSuccess)
	return promRegistry
}
