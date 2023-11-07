package service

import (
	"context"
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"

	"github.com/evgenytr/metrics.git/internal/interfaces"
	"github.com/evgenytr/metrics.git/internal/metric"
	pb "github.com/evgenytr/metrics.git/pkg/api/v1"
)

type MetricsServer struct {
	pb.UnimplementedMetricsServiceV1Server
	storage interfaces.Storage
}

// NewMetricsServerWithStorage returns MetricsServer for passed storage interface.
func NewMetricsServerWithStorage(storage interfaces.Storage) *MetricsServer {
	return &MetricsServer{
		storage: storage,
	}
}

func (s *MetricsServer) MetricsBatchV1(ctx context.Context, req *pb.MetricsBatchRequest) (resp *status.Status, err error) {
	fmt.Println("Metrics Batch V1 gRPC")
	fmt.Println(req)
	resp = &status.Status{
		Code: int32(codes.OK),
	}

	for _, currMetric := range req.GetMetrics() {
		switch currMetric.GetType() {
		case metric.GaugeMetricType:
			currMetric.Value, err = s.storage.UpdateGauge(ctx, currMetric.GetId(), currMetric.GetValue())
		case metric.CounterMetricType:
			currMetric.Delta, err = s.storage.UpdateCounter(ctx, currMetric.GetId(), currMetric.GetDelta())
		default:
			err = fmt.Errorf("metric type unknown %v", currMetric.GetType())
		}

		if err != nil {
			resp.Code = int32(codes.InvalidArgument)
			return
		}
	}

	return
}
