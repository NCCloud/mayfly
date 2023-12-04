package common

import (
	"fmt"
	"time"

	"github.com/NCCloud/mayfly/pkg/apis/v1alpha1"

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

func (s *Scheduler) CreateOrUpdateCreationJob(date time.Time,
	task func(resource v1alpha1.ScheduledResource) error, resource v1alpha1.ScheduledResource,
) error {
	tag := fmt.Sprintf("%v-create", resource.GetUID())

	if jobs, _ := s.scheduler.FindJobsByTag(tag); len(jobs) > 0 {
		if jobs[0].NextRun().Equal(date) {
			return nil
		}

		if removeJobErr := s.scheduler.RemoveByTag(tag); removeJobErr != nil {
			return removeJobErr
		}
	}

	_, jobErr := s.scheduler.StartAt(date).Every(1).LimitRunsTo(1).Tag(tag).
		Do(task, resource)

	return jobErr
}

func (s *Scheduler) CreateOrUpdateDeletionJob(date time.Time,
	task func(resource client.Object) error, resource client.Object,
) error {
	tag := fmt.Sprintf("%v-delete", resource.GetUID())

	if jobs, _ := s.scheduler.FindJobsByTag(tag); len(jobs) > 0 {
		if jobs[0].NextRun().Equal(date) {
			return nil
		}

		if removeJobErr := s.scheduler.RemoveByTag(tag); removeJobErr != nil {
			return removeJobErr
		}
	}

	_, jobErr := s.scheduler.StartAt(date).Every(1).LimitRunsTo(1).Tag(tag).
		Do(task, resource)

	return jobErr
}

func (s *Scheduler) DeleteCreationJob(resource client.Object) error {
	return s.scheduler.RemoveByTag(fmt.Sprintf("%v-create", resource.GetUID()))
}

func (s *Scheduler) DeleteDeletionJob(resource client.Object) error {
	return s.scheduler.RemoveByTag(fmt.Sprintf("%v-delete", resource.GetUID()))
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
