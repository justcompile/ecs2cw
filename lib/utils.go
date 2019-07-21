package lib

import (
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

func chunkSlice(s []*string, chunkSize int) [][]*string {
	var result [][]*string

	for i := 0; i < len(s); i += chunkSize {
		end := i + chunkSize

		if end > len(s) {
			end = len(s)
		}

		result = append(result, s[i:end])
	}

	return result
}

func chunkMetrics(s []*cloudwatch.MetricDatum, chunkSize int) [][]*cloudwatch.MetricDatum {
	var result [][]*cloudwatch.MetricDatum

	for i := 0; i < len(s); i += chunkSize {
		end := i + chunkSize

		if end > len(s) {
			end = len(s)
		}

		result = append(result, s[i:end])
	}

	return result
}
