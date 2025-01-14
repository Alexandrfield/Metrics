package databasestorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Alexandrfield/Metrics/internal/common"
)

type MemDatabaseStorage struct {
	Logger      common.Loger
	db          *sql.DB
	DatabaseDsn string
}

func (m *MemDatabaseStorage) createTable(ctx context.Context) error {
	const query = `CREATE TABLE if NOT EXISTS metrics (id text PRIMARY KEY, mtype text, delta bigint, value DOUBLE PRECISION)`
	if _, err := m.db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("error while trying to create table: %w", err)
	}
	return nil
}
func (st *MemDatabaseStorage) Start() error {
	var err error
	st.db, err = sql.Open("pgx", st.DatabaseDsn)
	if err != nil {
		st.db = nil
		return fmt.Errorf("can not open database. err:%w", err)
	}
	st.Logger.Infof("Connect to db open")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = st.createTable(ctx)
	if err != nil {
		st.db = nil
		return fmt.Errorf("can not create table. err:%w", err)
	}
	return nil
}
func (st *MemDatabaseStorage) AddGauge(name string, value common.TypeGauge) error {
	query := `INSERT INTO metrics (id, mtype, value) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET value = $4;`
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if _, err := st.db.ExecContext(ctx, query, name, "gauge", value, value); err != nil {
		return fmt.Errorf("error while trying to save gauge metric %s: %w", name, err)
	}
	return nil
}
func (st *MemDatabaseStorage) GetGauge(name string) (common.TypeGauge, error) {
	row := st.db.QueryRowContext(context.Background(),
		"SELECT value FROM metrics WHERE id = $1 AND mtype = gauge", name)
	var res common.TypeGauge
	err := row.Scan(&res)
	if err != nil {
		return common.TypeGauge(0), fmt.Errorf("Problem with scan GetGauge. err:%w", err)
	}
	return res, nil
}
func (st *MemDatabaseStorage) AddCounter(name string, value common.TypeCounter) error {
	query := `INSERT INTO metrics (id, mtype, delta) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET delta = EXCLUDED.delta + $4;`
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if _, err := st.db.ExecContext(ctx, query, name, "counter", value, value); err != nil {
		return fmt.Errorf("error while trying to save counter metric %s: %w", name, err)
	}
	return nil
}
func (st *MemDatabaseStorage) GetCounter(name string) (common.TypeCounter, error) {
	row := st.db.QueryRowContext(context.Background(),
		"SELECT delta FROM metrics WHERE id = $1 AND mtype = counter", name)
	var res common.TypeCounter
	err := row.Scan(&res)
	if err != nil {
		return common.TypeCounter(0), fmt.Errorf("Problem with scan GetGauge. err:%w", err)
	}
	return res, nil
}

func (st *MemDatabaseStorage) GetAllMetricName() ([]string, []string) {
	allGaugeKeys := make([]string, 0)
	allCounterKeys := make([]string, 0)
	rows, err := st.db.QueryContext(context.Background(), "SELECT id FROM metrics WHERE  mtype = gauge")
	if err != nil {
		st.Logger.Errorf("Problem with QueryContext gauge metrics. err:%w", err)
		return allGaugeKeys, allCounterKeys
	}
	for rows.Next() {
		var temp string
		err := rows.Scan(&temp)
		if err != nil {
			st.Logger.Errorf("Problem with scan gauge name", err)
			continue
		}
		allGaugeKeys = append(allGaugeKeys, temp)
	}
	rows, err = st.db.QueryContext(context.Background(), "SELECT id FROM metrics WHERE  mtype = counter")
	if err != nil {
		st.Logger.Errorf("Problem with QueryContext counter metrics. err:%w", err)
		return allGaugeKeys, allCounterKeys
	}
	for rows.Next() {
		var temp string
		err := rows.Scan(&temp)
		if err != nil {
			st.Logger.Errorf("Problem with scan counter name", err)
			continue
		}
		allGaugeKeys = append(allGaugeKeys, temp)
	}
	return allGaugeKeys, allCounterKeys
}

func (st *MemDatabaseStorage) PingDatabase() bool {
	status := false
	if st.db != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		if err := st.db.PingContext(ctx); err == nil {
			status = true
			st.Logger.Infof("succesful ping")
		} else {
			st.Logger.Infof("Can not ping database. err:%w", err)
		}
	}
	return status
}
