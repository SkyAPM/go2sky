package go2sky

import (
	"fmt"
	"testing"
)

func TestInitMetricCollector(t *testing.T) {
	mockMetricsReporter := MockMetricsReporter{}
	InitMetricCollector(&mockMetricsReporter)
}

type MockMetricsReporter struct {
}

func (m *MockMetricsReporter) SendMetrics(metrics RunTimeMetric) {
	fmt.Println(metrics)
}
