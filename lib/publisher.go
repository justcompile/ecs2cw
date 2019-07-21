package lib

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
)

type publisher struct {
	client    cloudwatchiface.CloudWatchAPI
	namespace *string
}

func (p *publisher) publishMetrics(account *account, metrics []*Metric) error {
	log.Printf("Account: %s, Region: %s => Publishing %d metrics", account.ID, account.Region, len(metrics))

	allMetrics := make([]*cloudwatch.MetricDatum, 0)
	for _, metric := range metrics {
		allMetrics = append(allMetrics, p.metricToCloudwatchMetrics(metric)...)
	}

	for _, metrics := range chunkMetrics(allMetrics, 20) {
		params := &cloudwatch.PutMetricDataInput{
			Namespace:  p.namespace,
			MetricData: metrics,
		}

		if _, err := p.client.PutMetricData(params); err != nil {
			return err
		}
	}

	return nil
}

func (p *publisher) metricToCloudwatchMetrics(metric *Metric) []*cloudwatch.MetricDatum {
	return []*cloudwatch.MetricDatum{
		{
			MetricName:        aws.String("DesiredCount"),
			Dimensions:        p.dimensionsForMetric(metric),
			Value:             aws.Float64(float64(*metric.DesiredCount)),
			StorageResolution: aws.Int64(60),
			Timestamp:         aws.Time(time.Now()),
		},
		{
			MetricName:        aws.String("PendingCount"),
			Dimensions:        p.dimensionsForMetric(metric),
			Value:             aws.Float64(float64(*metric.PendingCount)),
			StorageResolution: aws.Int64(60),
			Timestamp:         aws.Time(time.Now()),
		},
		{
			MetricName:        aws.String("RunningCount"),
			Dimensions:        p.dimensionsForMetric(metric),
			Value:             aws.Float64(float64(*metric.RunningCount)),
			StorageResolution: aws.Int64(60),
			Timestamp:         aws.Time(time.Now()),
		},
	}
}

func (p *publisher) dimensionsForMetric(metric *Metric) []*cloudwatch.Dimension {
	return []*cloudwatch.Dimension{
		{
			Name:  aws.String("Cluster"),
			Value: metric.Cluster,
		},
		{
			Name:  aws.String("Service"),
			Value: metric.ServiceName,
		},
	}
}

func newPublisher(namespace string, config *aws.Config, session client.ConfigProvider) *publisher {

	return &publisher{
		cloudwatch.New(session, config),
		aws.String(namespace),
	}
}
