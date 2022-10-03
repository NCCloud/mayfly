package scheduler

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var mayflyTotalJobs = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "mayfly_jobs_total",
		Help: "Number of scheduled Mayfly Jobs",
	},
)

var mayflyPastJobs = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "mayfly_jobs_total_past",
		Help: "Number of past Mayfly Jobs",
	},
)

func init() {
	metrics.Registry.MustRegister(mayflyTotalJobs, mayflyPastJobs)
}

func exportMayflyTotalJobsMetrics(jobs float64) {
	mayflyTotalJobs.Set(jobs)
}

func exportMayflyPastJobsMetrics(jobs float64) {
	mayflyPastJobs.Set(jobs)
}
