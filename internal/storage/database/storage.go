package databasestorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Alexandrfield/Metrics/internal/common"
)

const typecounter = "counter"
const typegauge = "gauge"

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
func (st *MemDatabaseStorage) exec(ctx context.Context, query string, args ...any) error {
	waitTime := []int{0, 1, 3, 5}
	for _, val := range waitTime {
		time.Sleep(time.Duration(val) * time.Second)
		st.Logger.Debugf("Retry query to database:%s;%s;", query, args)
		if _, err := st.db.ExecContext(ctx, query, args...); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
				st.Logger.Debugf("Retry query to database")
				continue
			}
			return fmt.Errorf("error exec database request. err: %w", err)
		} else {
			break
		}
	}
	return nil
}

func (st *MemDatabaseStorage) AddGauge(name string, value common.TypeGauge) error {
	_, err := st.GetGauge(name)
	if err != nil {
		st.Logger.Debugf("gauge metric new nmae:%s; value:%d;", name, value)
		query := "INSERT INTO metrics (id, mtype, value) VALUES ($1, $2, $3)"
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := st.exec(ctx, query, name, typegauge, value); err != nil {
			return fmt.Errorf("error while trying to save counter metric %s: %w", name, err)
		}
	} else {
		st.Logger.Debugf("gauge metric exist nmae:%s; value:%d; res:%d", name, value)
		query := "UPDATE metrics SET value = $1 WHERE id = $2"
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := st.exec(ctx, query, value, name); err != nil {
			return fmt.Errorf("error while trying to save counter metric %s: %w", name, err)
		}
	}
	return nil
}
func (st *MemDatabaseStorage) GetGauge(name string) (common.TypeGauge, error) {
	row := st.db.QueryRowContext(context.Background(),
		"SELECT value FROM metrics WHERE (id = $1 AND mtype = $2)", name, typegauge)
	st.Logger.Debugf("Gnry find etGauge name:%s", name)
	var res common.TypeGauge
	err := row.Scan(&res)
	if err != nil {
		return common.TypeGauge(0), fmt.Errorf("problem with scan GetGauge. err:%w", err)
	}
	st.Logger.Debugf("GetGauge name:%s->:%d", name, res)
	return res, nil
}
func (st *MemDatabaseStorage) AddCounter(name string, value common.TypeCounter) error {
	st.Logger.Debugf("addCounter to database. name:%s, value:%d", name, value)
	val, err := st.GetCounter(name)
	if err != nil {
		st.Logger.Debugf("counter metric new nmae:%s; value:%d;", name, value)
		query := "INSERT INTO metrics (id, mtype, delta) VALUES ($1, $2, $3)"
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := st.exec(ctx, query, name, typecounter, value); err != nil {
			return fmt.Errorf("error while trying to save counter metric %s: %w", name, err)
		}
	} else {
		st.Logger.Debugf("counter metric exist nmae:%s; val:%d; value:%d; res:%d", name, val, value, val+value)
		query := "UPDATE metrics SET delta = $1 WHERE id = $2"
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := st.exec(ctx, query, val+value, name); err != nil {
			return fmt.Errorf("error while trying to save counter metric %s: %w", name, err)
		}
	}
	st.Logger.Debugf("succes add counter metric name%s;", name)
	return nil
}
func (st *MemDatabaseStorage) GetCounter(name string) (common.TypeCounter, error) {
	row := st.db.QueryRowContext(context.Background(),
		"SELECT delta FROM metrics WHERE (id = $1 AND mtype = $2)", name, typecounter)
	var res common.TypeCounter
	st.Logger.Debugf("row:%s", row)
	err := row.Scan(&res)
	if err != nil {
		st.Logger.Debugf("cant find name:%s; err:%w", name, err)
		return common.TypeCounter(0), fmt.Errorf("problem with scan GetCounter. err:%w", err)
	}
	st.Logger.Debugf("GetCounter name:%s->:%d", name, res)
	return res, nil
}

func (st *MemDatabaseStorage) GetAllMetricName() ([]string, []string) {
	allGaugeKeys := make([]string, 0)
	allCounterKeys := make([]string, 0)
	rowsg, err := st.db.QueryContext(context.Background(), "SELECT id FROM metrics WHERE  mtype = gauge")
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
			continue
		case "gauge":
			query := `INSERT INTO metrics (id, mtype, value) VALUES ($1, $2, $3) ON CONFLICT (id) 
			DO UPDATE SET value = $4;`
			if err := st.exec(context.Background(), query, metric.ID, typegauge, common.TypeGauge(*metric.Value),
				common.TypeGauge(*metric.Value)); err != nil {
				errr := tx.Rollback()
				if errr != nil {
					return fmt.Errorf("error while trying to save gauge metric %s: err%w;And can not rollback! err:%w",
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
	for _, metric := range metrics {
		switch metric.MType {
		case "counter":
			_ = st.AddCounter(metric.ID, common.TypeCounter(*metric.Delta))
		default:
			continue
		}
	}

	return nil
}
