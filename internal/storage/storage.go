package storage

import (
	"os"

	"github.com/Alexandrfield/Metrics/internal/common"
	database_storage "github.com/Alexandrfield/Metrics/internal/storage/database"
	file_storage "github.com/Alexandrfield/Metrics/internal/storage/file"
)

//go:generate mockgen -source=storage.go -destination=mock/storage.go

// BasicStorage is common interface for object save metric.
type BasicStorage interface {
	AddCounter(metricName string, metricValue common.TypeCounter) error
	AddGauge(name string, value common.TypeGauge) error
	GetCounter(metricName string) (common.TypeCounter, error)
	GetGauge(metricName string) (common.TypeGauge, error)
	GetAllMetricName() ([]string, []string)
	PingDatabase() bool
	AddMetrics(metrics []common.Metrics) error
	Close()
}

// CreateMemStorage create database, file  storage depending on the configuration parameters.
func CreateMemStorage(config Config, logger common.Loger, done chan struct{}) BasicStorage {
	var res BasicStorage
	if config.DatabaseDsn != "" {
		logger.Debugf("Create storage database")
		memStorage := database_storage.NewMemDatabaseStorage(logger, config.DatabaseDsn)
		err := memStorage.Start()
		if err != nil {
			logger.Debugf("Issue with start database %s", err)
		}
		res = memStorage
	} else {
		logger.Debugf("Create storage file")
		memStorage := file_storage.NewMemFileStorage(config.FileStoregePath, logger)
		logger.Debugf("config.Restore %s", config.Restore)
		if config.Restore {
			file, err := os.OpenFile(config.FileStoregePath, os.O_RDONLY, 0o600)
			if err == nil {
				memStorage.Logger.Debugf("file was open")
				defer func() {
					_ = file.Close()
				}()
				memStorage.LoadMemStorage(file)
			} else {
				logger.Debugf("can not restore file. File is not exist. err:%w", err)
			}
		}
		go file_storage.StorageSaver(memStorage, config.StoreIntervalSecond, done)
		res = memStorage
	}
	go func() {
		<-done
		res.Close()
	}()
	return res
}
