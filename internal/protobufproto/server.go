package protobufproto

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/Alexandrfield/Metrics/internal/common"
	pr "github.com/Alexandrfield/Metrics/internal/protobufproto/proto/protobufproto"
	"google.golang.org/grpc"
)

type MetricsStorage interface {
	SetCounterValue(metricName string, metricValue common.TypeCounter) error
	SetGaugeValue(metricName string, metricValue common.TypeGauge) error
	GetCounterValue(metricName string) (common.TypeCounter, error)
	GetGaugeValue(metricName string) (common.TypeGauge, error)
	GetAllValue() ([]string, error)
	PingDatabase() bool
	AddMetrics(metrics []common.Metrics) error
}

type GRPCServer struct {
	logger common.Loger
	pr.GRPCMetricsServiceServer
	storage MetricsStorage
}

func (g *GRPCServer) Update(ctx context.Context, req *pr.RequestMetric) (*pr.Status, error) {
	g.logger.Debugf("Update req.MType:%s, req.ID:%s;", req.MType, req.ID)
	if req.MType == "Gauge" {
		err := g.storage.SetGaugeValue(req.ID, common.TypeGauge(req.Value))
		if err != nil {
			return &pr.Status{Status: http.StatusInternalServerError},
				fmt.Errorf("server error update gauge. err:%s", err)
		}
	}
	if req.MType == "Counter" {
		err := g.storage.SetCounterValue(req.ID, common.TypeCounter(req.Delta))
		if err != nil {
			return &pr.Status{Status: http.StatusInternalServerError},
				fmt.Errorf("server error update counter. err:%s", err)
		}
	}
	return &pr.Status{Status: http.StatusOK}, nil
}

func (g *GRPCServer) Value(ctx context.Context, req *pr.RequestMetric) (*pr.RequestAnsMetric, error) {
	sendMetric := pr.RequestAnsMetric{ID: req.ID, MType: req.MType}
	if req.MType == "Gauge" {
		data, err := g.storage.GetGaugeValue(req.ID)
		if err != nil {
			sendMetric.Status = http.StatusInternalServerError
			return &sendMetric, fmt.Errorf("server error update gauge. err:%s", err)
		}
		sendMetric.Value = float64(data)
	}
	if req.MType == "Counter" {
		data, err := g.storage.GetCounterValue(req.ID)
		if err != nil {
			sendMetric.Status = http.StatusInternalServerError
			return &sendMetric, fmt.Errorf("server error update counter. err:%s", err)
		}
		sendMetric.Delta = int64(data)
	}
	sendMetric.Status = http.StatusOK
	return &sendMetric, nil
}

func StartGRPCServer(logger common.Loger, str MetricsStorage) {
	grpcWorker := GRPCServer{logger: logger, storage: str}
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Fatalf("Error start tcp server for grpc: %v", err)
	}

	grpcServer := grpc.NewServer()
	pr.RegisterGRPCMetricsServiceServer(grpcServer, &grpcWorker)

	logger.Infof("start grpc server.")
	if err := grpcServer.Serve(listener); err != nil {
		logger.Fatalf("Error start tcp server for grpc: %v", err)
	}
}
