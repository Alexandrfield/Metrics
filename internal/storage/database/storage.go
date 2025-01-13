package databasestorage

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Alexandrfield/Metrics/internal/common"
)

type MemDatabaseStorage struct {
	Logger      common.Loger
	db          *sql.DB
	DatabaseDsn string
}

func (st *MemDatabaseStorage) Start() error {
	// ps := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
	// 	`localhost`, `video`, `XXXXXXXX`, `video`)
	var err error
	st.db, err = sql.Open("pgx", st.DatabaseDsn)
	if err != nil {
		st.db = nil
		return fmt.Errorf("can not open database. err:%w", err)
	}
	st.Logger.Infof("Connect to db open")
	// defer st.db.Close()
	return nil
}
func (st *MemDatabaseStorage) AddGauge(name string, value common.TypeGauge) error {
	return nil
}
func (st *MemDatabaseStorage) GetGauge(name string) (common.TypeGauge, error) {
	return common.TypeGauge(0), nil
}
func (st *MemDatabaseStorage) AddCounter(name string, value common.TypeCounter) error {
	return nil
}
func (st *MemDatabaseStorage) GetCounter(name string) (common.TypeCounter, error) {
	return common.TypeCounter(0), nil
}

func (st *MemDatabaseStorage) GetAllMetricName() ([]string, []string) {
	allGaugeKeys := make([]string, 0)
	allCounterKeys := make([]string, 0)
	return allGaugeKeys, allCounterKeys
}

func (st *MemDatabaseStorage) PingDatabase() bool {
	return st.db != nil
}
