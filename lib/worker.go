package lib

import (
	"log"
	"os"
)

type Worker struct {
	account   *account
	namespace string
}

func (w *Worker) do() error {
	g := newGatherer(w.account)

	clusters, err := g.getClusterARNs()
	if err != nil {
		return err
	}

	allMetrics := make([]*Metric, 0)

	for _, cluster := range clusters {
		services, err := g.getServicesForCluster(cluster)
		if err != nil {
			return err
		}

		metrics, err := g.getServiceMetrics(cluster, services)
		if err != nil {
			return err
		}

		allMetrics = append(allMetrics, metrics...)
	}

	p := newPublisher(w.namespace, g.config, g.session)
	return p.publishMetrics(w.account, allMetrics)
}

func newWorker(account *account) *Worker {
	namespace := os.Getenv("METRIC_NAMESPACE")
	if namespace == "" {
		log.Fatalf("METRIC_NAMESPACE environment variable is not defined")
	}

	return &Worker{
		account,
		namespace,
	}
}
