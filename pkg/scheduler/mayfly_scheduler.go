package scheduler

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-co-op/gocron"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/NCCloud/mayfly/pkg/common"
)

type Scheduler struct {
	Config    *common.OperatorConfig
	Client    client.Client
	Scheduler *gocron.Scheduler
}

func NewScheduler(config *common.OperatorConfig, client client.Client) *Scheduler {
	s := gocron.NewScheduler(time.UTC)
	s.StartAsync()

	return &Scheduler{
		Config:    config,
		Client:    client,
		Scheduler: s,
	}
}

func (s *Scheduler) StartMonitor() {
	go func() {
		for {
			//TODO: Expose as metrics
			fmt.Println("=============")
			fmt.Println("Total jobs: " + strconv.FormatInt(int64(len(s.Scheduler.Jobs())), 10))
			for _, job := range s.Scheduler.Jobs() {
				fmt.Println("\tNext run: " + job.NextRun().String())
			}
			time.Sleep(time.Second * 5)
		}
	}()
}

func (s *Scheduler) StartOrUpdateJob(expirationDate time.Time, task func(ctx context.Context, client client.Client, resource client.Object) error, ctx context.Context, client client.Client, resource client.Object) error {
	jobs, _ := s.Scheduler.FindJobsByTag(fmt.Sprintf("%v", resource.GetUID()))

	if len(jobs) > 0 {
		jobExpirationDate := jobs[0].NextRun()
		if jobExpirationDate.Equal(expirationDate) {
			return nil
		}
		s.RemoveJob(fmt.Sprintf("%v", resource.GetUID()))
	}

	_, jobErr := s.Scheduler.
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

func (s *Scheduler) RemoveJob(id string) {
	_ = s.Scheduler.RemoveByTag(id)
}
