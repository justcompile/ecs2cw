package lib

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
)

type gatherer struct {
	client  ecsiface.ECSAPI
	session *session.Session
	config  *aws.Config
}

func (g *gatherer) getServiceMetrics(cluster *string, serviceARNs []*string) ([]*Metric, error) {
	metrics := make([]*Metric, 0)

	for _, arns := range chunkSlice(serviceARNs, 10) {
		params := &ecs.DescribeServicesInput{
			Cluster:  cluster,
			Services: arns,
		}

		output, err := g.client.DescribeServices(params)
		if err != nil {
			return nil, err
		}

		metrics = append(metrics, MetricsFromServices(cluster, output.Services)...)
	}

	return metrics, nil
}

func (g *gatherer) getServicesForCluster(arn *string) ([]*string, error) {
	params := &ecs.ListServicesInput{
		Cluster: arn,
	}

	services := make([]*string, 0)

	err := g.client.ListServicesPages(
		params,
		func(output *ecs.ListServicesOutput, lastPage bool) bool {
			services = append(services, output.ServiceArns...)

			return lastPage
		},
	)

	return services, err
}

func (g *gatherer) getClusterARNs() ([]*string, error) {
	listParams := new(ecs.ListClustersInput)

	output, err := g.client.ListClusters(listParams)
	if err != nil {
		return nil, err
	}

	return output.ClusterArns, nil
}

func newGatherer(account *account) *gatherer {
	sess := session.Must(session.NewSession())

	cfg := &aws.Config{Region: aws.String(account.Region)}

	if role := os.Getenv("ROLE_OVERRIDE"); role != "" {
		log.Printf("Assuming role: %s\n", role)
		creds := stscreds.NewCredentials(sess, role, func(p *stscreds.AssumeRoleProvider) {
			p.TokenProvider = stscreds.StdinTokenProvider
		})

		cfg.Credentials = creds
	} else {
		creds := stscreds.NewCredentials(sess, fmt.Sprintf("arn:aws:iam::%s:role/VanguardWorkerRole", account.ID), func(p *stscreds.AssumeRoleProvider) {
			p.TokenProvider = stscreds.StdinTokenProvider
		})

		cfg.Credentials = creds
	}

	log.Printf("Account: %s, Region: %s => Gathering metrics", account.ID, account.Region)

	return &gatherer{
		client:  ecs.New(sess, cfg),
		session: sess,
		config:  cfg,
	}
}
