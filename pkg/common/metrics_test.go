package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

func TestMetrics_Init(t *testing.T) {
	// when
	gather, gatherErr := metrics.Registry.Gather()

	// then
	assert.Nil(t, gatherErr)
	assert.Equal(t, "mayfly_total_jobs", gather[0].GetName())
}
