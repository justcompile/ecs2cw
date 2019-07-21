package lib

import "github.com/aws/aws-sdk-go/service/ecs"

type Metric struct {
	Cluster      *string
	ServiceName  *string
	DesiredCount *int64
	PendingCount *int64
	RunningCount *int64
}

func MetricsFromServices(cluster *string, services []*ecs.Service) []*Metric {
	metrics := make([]*Metric, len(services))

	for i, service := range services {
		metrics[i] = &Metric{
			Cluster:      cluster,
			ServiceName:  service.ServiceName,
			DesiredCount: service.DesiredCount,
			PendingCount: service.PendingCount,
			RunningCount: service.RunningCount,
		}
	}

	return metrics
}
