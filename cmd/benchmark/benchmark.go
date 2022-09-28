package main

import (
	"context"
	"fmt"
	time "time"

	"github.com/NCCloud/mayfly/pkg"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var result = Result{}

func (b *Benchmark) Start() {
	b.startedAt = time.Now()
	go b.createSecrets()
	go b.Listener()
}

func (b *Benchmark) CreateResources() {
	fmt.Println("Creating resources")
}

func (b *Benchmark) Listener() Result {
	result = Result{}

	for {
		kind := make(map[string]int)
		now := time.Now()
		for _, resourceKind := range b.config.Resources {
			resourceList := pkg.NewResourceInstanceList(resourceKind)

			resourcesListErr := b.mgrClient.List(context.Background(), resourceList)
			if resourcesListErr != nil {
				panic(resourcesListErr)
			}
			for _, resource := range resourceList.Items {
				if resource.GetAnnotations()[b.config.ExpirationLabel] == "" {
					continue
				}
				kind[resourceKind]++
			}

		}
		result.Points = append(result.Points, Point{
			time: now,
			kind: kind,
		})

		result.Duration = time.Since(b.startedAt)
		result.StartedAt = b.startedAt

		if time.Since(b.startedAt) > 60*time.Minute {
			break
		}
		time.Sleep(b.granularity)
	}

	return result
}

func (b *Benchmark) createSecrets() {
	for i := 0; i < b.count; i++ {
		id := i + b.offset
		secret := b.generateSecret(id)
		createErr := mgrClient.Create(context.Background(), secret)
		if createErr != nil {
			fmt.Println(createErr)
		}

		fmt.Printf("\nSecret %d has been created.", id)
		time.Sleep(b.delay)
	}
	fmt.Printf("\n")
}

func (b *Benchmark) generateSecret(id int) *v1.Secret {
	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("benchmark-%d", id),
			Namespace: "default",
			Annotations: map[string]string{
				b.config.ExpirationLabel: "10s",
			},
		},
		Data: map[string][]byte{
			"password": []byte("password"),
		},
	}
}

func (b *Benchmark) GetResult() Result {
	return result
}
