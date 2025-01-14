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

func (st *MemDatabaseStorage) createTable(ctx context.Context) error {
	const query = `CREATE TABLE if NOT EXISTS metrics (id text PRIMARY KEY, 
	mtype text, delta bigint, value DOUBLE PRECISION)`
	if _, err := st.db.ExecContext(ctx, query); err != nil {
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
		return common.TypeGauge(0), fmt.Errorf("problem with scan GetGauge. err:%w", err)
	}
	return res, nil
}
func (st *MemDatabaseStorage) AddCounter(name string, value common.TypeCounter) error {
	query := `INSERT INTO metrics (id, mtype, delta) VALUES ($1, $2, $3) 
	ON CONFLICT (id) DO UPDATE SET delta = EXCLUDED.delta + $4;`
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if _, err := st.db.ExecContext(ctx, query, name, "counter", value, value); err != nil {
		st.Logger.Debugf("error while trying to save counter metric %s: %w", name, err)
		return fmt.Errorf("error while trying to save counter metric %s: %w", name, err)
	}
	st.Logger.Debugf("succes add counter metric name%s;", name)
	return nil
}
func (st *MemDatabaseStorage) GetCounter(name string) (common.TypeCounter, error) {
	row := st.db.QueryRowContext(context.Background(),
		"SELECT delta FROM metrics WHERE id = $1 AND mtype = \"counter\"", name)
	var res common.TypeCounter
	st.Logger.Debugf("row:%s", row)
	err := row.Scan(&res)
	if err != nil {
		st.Logger.Debugf("cant find name:%s; err:%w", name, err)
		return common.TypeCounter(0), fmt.Errorf("problem with scan GetCounter. err:%w", err)
	}
	return res, nil
}

func (st *MemDatabaseStorage) GetAllMetricName() ([]string, []string) {
	allGaugeKeys := make([]string, 0)
	allCounterKeys := make([]string, 0)
	rowsg, err := st.db.QueryContext(context.Background(), "SELECT id FROM metrics WHERE  mtype = \"gauge\"")
	defer func() { _ = rowsg.Close() }()
	if err != nil {
		st.Logger.Errorf("Problem with QueryContext gauge metrics. err:%w", err)
		return allGaugeKeys, allCounterKeys
	}
	for rowsg.Next() {
		var temp string
		err := rowsg.Scan(&temp)
		if err != nil {
			st.Logger.Errorf("Problem with scan gauge name", err)
			continue
		}
		allGaugeKeys = append(allGaugeKeys, temp)
	}
	if err = rowsg.Err(); err != nil {
		st.Logger.Errorf("Problem iteration gauge names", err)
	}
	rowsc, err := st.db.QueryContext(context.Background(), "SELECT id FROM metrics WHERE  mtype = counter")
	defer func() { _ = rowsc.Close() }()
	if err != nil {
		st.Logger.Errorf("Problem with QueryContext counter metrics. err:%w", err)
		return allGaugeKeys, allCounterKeys
	}
	for rowsc.Next() {
		var temp string
		err := rowsc.Scan(&temp)
		if err != nil {
			st.Logger.Errorf("problem with scan counter name", err)
			continue
		}
		allGaugeKeys = append(allGaugeKeys, temp)
	}
	if err = rowsc.Err(); err != nil {
		st.Logger.Errorf("Problem iteration counter names", err)
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

func (st *MemDatabaseStorage) AddMetrics(metrics []common.Metrics) error {
	tx, err := st.db.Begin()
	if err != nil {
		return fmt.Errorf("can not start transactiom. err:%w", err)
	}

	for _, metric := range metrics {
		switch metric.MType {
		case "counter":
			query := `INSERT INTO metrics (id, mtype, delta) VALUES ($1, $2, $3) 
	ON CONFLICT (id) DO UPDATE SET delta = EXCLUDED.delta + $4;`
			if _, err := tx.ExecContext(context.Background(), query, metric.ID, "counter", common.TypeCounter(*metric.Delta),
				common.TypeCounter(*metric.Delta)); err != nil {
				errr := tx.Rollback()
				if errr != nil {
					return fmt.Errorf("error while trying to save counter metric %s: err%w; And can not rollback! err:%w",
						metric.ID, err, errr)
				}
				return fmt.Errorf("error while trying to save counter metric %s: %w", metric.ID, err)
			}
		case "gauge":
			query := `INSERT INTO metrics (id, mtype, value) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET value = $4;`
			if _, err := st.db.ExecContext(context.Background(), query, metric.ID, "gauge", common.TypeGauge(*metric.Value),
				common.TypeGauge(*metric.Value)); err != nil {
				errr := tx.Rollback()
				if errr != nil {
					return fmt.Errorf("error while trying to save gauge metric %s: err%w; And can not rollback! err:%w",
						metric.ID, err, errr)
				}
				return fmt.Errorf("error while trying to save gauge metric %s: %w", metric.ID, err)
			}
		default:
			st.Logger.Debugf("AddMetrics. unknown type:%s;", metric.MType)
			continue
		}
	}
	comerr := tx.Commit()
	if comerr != nil {
		return fmt.Errorf("error with commit transactiom. err:%w", err)
	}
	return nil
}
