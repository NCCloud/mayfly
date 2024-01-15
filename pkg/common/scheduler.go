package common

import (
	"time"

	"github.com/go-co-op/gocron"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Scheduler struct {
	config    *Config
	client    client.Client
	scheduler *gocron.Scheduler
}

func NewScheduler(config *Config, client client.Client) *Scheduler {
	scheduler := &Scheduler{
		config:    config,
		client:    client,
		scheduler: gocron.NewScheduler(time.UTC),
	}

	scheduler.startMonitoring()
	scheduler.scheduler.StartAsync()

	return scheduler
}

func (s *Scheduler) CreateOrUpdateTask(tag string, date time.Time, task func() error) error {
	if jobs, _ := s.scheduler.FindJobsByTag(tag); len(jobs) > 0 {
		if jobs[0].NextRun().Equal(date) {
			return nil
		}

		if removeJobErr := s.scheduler.RemoveByTag(tag); removeJobErr != nil {
			return removeJobErr
		}
	}

	_, jobErr := s.scheduler.StartAt(date).Every(1).LimitRunsTo(1).Tag(tag).Do(task)

	return jobErr
}

func (s *Scheduler) DeleteTask(tag string) error {
	return s.scheduler.RemoveByTag(tag)
}

func (s *Scheduler) startMonitoring() {
	if _, doErr := s.scheduler.SingletonMode().Every(s.config.MonitoringInterval).Do(func() {
		exportMayflyTotalJobsMetrics(float64(len(s.scheduler.Jobs())))

		pastJobs := 0
		for _, job := range s.scheduler.Jobs() {
			if job.NextRun().Before(time.Now()) {
				pastJobs++
			}
		}

		exportMayflyPastJobsMetrics(float64(pastJobs))
	}); doErr != nil {
		panic(doErr)
	}
}
