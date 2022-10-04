package main

import (
	"time"

	"github.com/NCCloud/mayfly/pkg/common"
	"github.com/NCCloud/mayfly/pkg/controllers/mayfly/resource"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Benchmark struct {
	granularity    time.Duration
	config         *resource.Resources
	mgrClient      client.Client
	operatorConfig *common.OperatorConfig
	startedAt      time.Time
	count          int
	offset         int
	delay          time.Duration
}

type Result struct {
	StartedAt time.Time
	Duration  time.Duration
	Points    []Point
}

type Point struct {
	kind map[string]int
	time time.Time
}

func NewBenchmark(mgrClient client.Client, config *resource.Resources, operatorConfig *common.OperatorConfig, count int) *Benchmark {
	return &Benchmark{
		granularity:    5 * time.Second,
		config:         config,
		mgrClient:      mgrClient,
		operatorConfig: operatorConfig,
		count:          count,
		offset:         0,
		delay:          0,
	}
}

func (b *Benchmark) Granularity(granularity time.Duration) *Benchmark {
	b.granularity = granularity
	return b
}

func (b *Benchmark) Offset(offset int) *Benchmark {
	b.offset = offset
	return b
}

func (b *Benchmark) Delay(delay time.Duration) *Benchmark {
	b.delay = delay
	return b
}
