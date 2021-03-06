package metrics

import (
	"time"

	"github.com/RichardKnop/jsonhal"
	"github.com/RichardKnop/pinglist-api/util"
)

// MetricResponse ...
type MetricResponse struct {
	jsonhal.Hal
	Timestamp string `json:"timestamp"`
	Value     int64  `json:"value"`
}

// NewMetricResponse creates new MetricResponse instance
func NewMetricResponse(timestamp time.Time, value int64) (*MetricResponse, error) {
	return &MetricResponse{
		Timestamp: util.FormatTime(timestamp),
		Value:     value,
	}, nil
}
