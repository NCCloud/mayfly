package pkg

import (
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Scheduler struct {
	Config        *Config
	Client        client.Client
	CronScheduler *gocron.Scheduler
}

func NewScheduler(config *Config, client client.Client) *Scheduler {
	cronScheduler := gocron.NewScheduler(time.UTC)
	cronScheduler.StartAsync()

	scheduler := &Scheduler{
		Config:        config,
		Client:        client,
		CronScheduler: cronScheduler,
	}

	scheduler.startMonitoring()

	return scheduler
}

func (s *Scheduler) startMonitoring() {
	const granularity = 5 * time.Second

	go func() {
		for {
			exportMayflyTotalJobsMetrics(float64(len(s.CronScheduler.Jobs())))

			pastJobs := 0

			for _, job := range s.CronScheduler.Jobs() {
				if job.NextRun().Before(time.Now()) {
					pastJobs++
				}
			}

			exportMayflyPastJobsMetrics(float64(pastJobs))
			time.Sleep(granularity)
		}
	}()
}

func (s *Scheduler) StartOrUpdateJob(ctx context.Context, expirationDate time.Time,
	task func(ctx context.Context, client client.Client, resource client.Object) error,
	client client.Client, resource client.Object,
) error {
	jobs, _ := s.CronScheduler.FindJobsByTag(fmt.Sprintf("%v", resource.GetUID()))

	if len(jobs) > 0 {
		jobExpirationDate := jobs[0].NextRun()
		if jobExpirationDate.Equal(expirationDate) {
			return nil
		}

		_ = s.RemoveJob(fmt.Sprintf("%v", resource.GetUID()))
	}

	_, jobErr := s.CronScheduler.
		Every(1).
		LimitRunsTo(1).
		StartAt(expirationDate).
		Tag(fmt.Sprintf("%v", resource.GetUID())).
		Do(task, ctx, client, resource)

	if jobErr != nil {
		return jobErr
	}

	return nil
}

func (s *Scheduler) RemoveJob(id string) error {
	return s.CronScheduler.RemoveByTag(id)
}
