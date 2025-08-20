package main

import (
	"time"

	"github.com/NCCloud/mayfly/pkg/common"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Benchmark struct {
	granularity time.Duration
	config      *common.Config
	mgrClient   client.Client
	startedAt   time.Time
	count       int
	offset      int
	delay       time.Duration
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

func NewBenchmark(mgrClient client.Client, config *common.Config, count int) *Benchmark {
	granularityInSeconds := 5

	return &Benchmark{
		granularity: time.Duration(granularityInSeconds) * time.Second,
		config:      config,
		mgrClient:   mgrClient,
		count:       count,
		offset:      0,
		delay:       0,
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
