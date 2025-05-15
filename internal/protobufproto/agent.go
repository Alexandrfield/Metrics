package protobufproto

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Alexandrfield/Metrics/internal/common"
	pr "github.com/Alexandrfield/Metrics/internal/protobufproto/proto/protobufproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func UpdateMetricGRPC(logger common.Loger, metr common.Metrics) error {
	logger.Debugf(" send metric with grpc")
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Errorf("Error: %v", err)
	}
	defer conn.Close()
	client := pr.NewGRPCMetricsServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	sendMetric := pr.RequestMetric{ID: metr.ID, MType: metr.MType}
	if metr.Delta != nil {
		sendMetric.Delta = *metr.Delta
	}
	if metr.Value != nil {
		sendMetric.Value = *metr.Value
	}
	res, err := client.Update(ctx, &sendMetric)
	if err != nil {
		logger.Errorf("error update: %v", err)
	}
	if res.Status != http.StatusOK {
		return fmt.Errorf("error send metric. status:%d", res.Status)
	}
	return nil
}
