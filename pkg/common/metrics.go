package common

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var mayflyTotalJobs = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "mayfly_total_jobs",
		Help: "Number of scheduled Mayfly Jobs",
	},
)

func init() {
	metrics.Registry.MustRegister(mayflyTotalJobs)
}
