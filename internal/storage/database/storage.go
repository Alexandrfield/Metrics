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
		return fmt.Errorf("can not open database. err:%w", err)
	}
	st.Logger.Infof("Connect to db open")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = st.createTable(ctx)
	if err != nil {
		errClose := st.db.Close()
		if errClose != nil {
			return fmt.Errorf("can not create table err:%w; end close connection to database err:%w", err, errClose)
		}
		return fmt.Errorf("can not create table. err:%w", err)
	}
	return nil
}

type databaseDB interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func (st *MemDatabaseStorage) exec(ctx context.Context, con databaseDB, query string, args ...any) error {
	waitTime := []int{0, 1, 3, 5}
	for _, val := range waitTime {
		time.Sleep(time.Duration(val) * time.Second)
		st.Logger.Debugf("Retry query to database:%s;%s;", query, args)
		if _, err := con.ExecContext(ctx, query, args...); err != nil {
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
	tx, err := st.db.Begin()
	if err != nil {
		return fmt.Errorf("can not create transaction. err:%w", err)
	}
	_, err = st.getGauge(tx, name)
	if err != nil {
		st.Logger.Debugf("gauge metric new nmae:%s; value:%d;", name, value)
		query := "INSERT INTO metrics (id, mtype, value) VALUES ($1, $2, $3)"
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := st.exec(ctx, tx, query, name, typegauge, value); err != nil {
			errRol := tx.Rollback()
			if errRol != nil {
				return fmt.Errorf("error while trying to save counter metric %s: %w; and error rollback err:%w", name, err, errRol)
			}
			return fmt.Errorf("error while trying to save counter metric %s: %w", name, err)
		}
	} else {
		st.Logger.Debugf("gauge metric exist nmae:%s; value:%d; res:%d", name, value)
		query := "UPDATE metrics SET value = $1 WHERE id = $2"
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := st.exec(ctx, tx, query, value, name); err != nil {
			errRol := tx.Rollback()
			if errRol != nil {
				return fmt.Errorf("error while trying to save counter metric %s: %w; and error rollback err:%w", name, err, errRol)
			}
			return fmt.Errorf("error while trying to save counter metric %s: %w", name, err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error with commit transactiom AddGauge. err:%w", err)
	}
	return nil
}

func (st *MemDatabaseStorage) getGauge(con databaseDB, name string) (common.TypeGauge, error) {
	row := con.QueryRowContext(context.Background(),
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

func (st *MemDatabaseStorage) GetGauge(name string) (common.TypeGauge, error) {
	res, err := st.getGauge(st.db, name)
	if err != nil {
		return common.TypeGauge(0), fmt.Errorf("GetGage err:%w", err)
	}
	return res, nil
}
func (st *MemDatabaseStorage) AddCounter(name string, value common.TypeCounter) error {
	st.Logger.Debugf("addCounter to database. name:%s, value:%d", name, value)
	tx, err := st.db.Begin()
	if err != nil {
		return fmt.Errorf("can not create transaction AddCounter. err:%w", err)
	}
	val, err := st.getCounter(tx, name)
	if err != nil {
		st.Logger.Debugf("counter metric new nmae:%s; value:%d;", name, value)
		query := "INSERT INTO metrics (id, mtype, delta) VALUES ($1, $2, $3)"
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := st.exec(ctx, tx, query, name, typecounter, value); err != nil {
			return fmt.Errorf("tx, error while trying to save counter metric %s: %w", name, err)
		}
	} else {
		st.Logger.Debugf("counter metric exist nmae:%s; val:%d; value:%d; res:%d", name, val, value, val+value)
		query := "UPDATE metrics SET delta = $1 WHERE id = $2"
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := st.exec(ctx, tx, query, val+value, name); err != nil {
			return fmt.Errorf("error while trying to save counter metric %s: %w", name, err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error with commit transactiom AddCounter. err:%w", err)
	}
	st.Logger.Debugf("succes add counter metric name%s;", name)
	return nil
}

func (st *MemDatabaseStorage) getCounter(con databaseDB, name string) (common.TypeCounter, error) {
	row := con.QueryRowContext(context.Background(),
		"SELECT delta FROM metrics WHERE (id = $1 AND mtype = $2)", name, typecounter)
	var res common.TypeCounter
	st.Logger.Debugf("SELECT delta FROM metrics WHERE (id = %s AND mtype = %s)", name, typecounter)
	err := row.Scan(&res)
	if err != nil {
		st.Logger.Debugf("cant find name:%s; err:%w", name, err)
		return common.TypeCounter(0), fmt.Errorf("problem with scan GetCounter. err:%w", err)
	}
	st.Logger.Debugf("GetCounter name:%s->:%d", name, res)
	return res, nil
}

func (st *MemDatabaseStorage) GetCounter(name string) (common.TypeCounter, error) {
	res, err := st.getCounter(st.db, name)
	if err != nil {
		return common.TypeCounter(0), fmt.Errorf("GetCounter err:%w", err)
	}
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := st.db.PingContext(ctx); err == nil {
		status = true
		st.Logger.Infof("succesful ping")
	} else {
		st.Logger.Infof("Can not ping database. err:%w", err)
	}
	return status
}

func (st *MemDatabaseStorage) AddMetrics(metrics []common.Metrics) error {
	tx, err := st.db.Begin()
	if err != nil {
		return fmt.Errorf("can not start transactiom. err:%w", err)
	}
	metricsGauge := make(map[string]string)
	metricsCounter := make(map[string]common.Metrics)
	for indexNum, metric := range metrics {
		st.Logger.Infof("%d) try add ID:%s; MType:%s; value:%d; delta:%d",
			indexNum, metric.ID, metric.MType, *metric.Value, *metric.Delta)
		switch metric.MType {
		case "counter":
			newVal := *metric.Delta
			val, ok := metricsCounter[metric.ID]
			if ok {
				newVal += *val.Delta
			}
			newVal += *metric.Delta
			metricsCounter[metric.ID] = common.Metrics{ID: metric.ID, MType: typecounter, Delta: &newVal}
		case "gauge":
			metricsGauge[metric.ID] = metric.GetValueMetric()
		default:
			continue
		}
	}
	metricsGaugeUpdate := make(map[string]string)
	metricsGaugeInsert := make(map[string]string)
	for key, val := range metricsGauge {
		_, err := st.getGauge(tx, key)
		if err != nil {
			metricsGaugeInsert[key] = val
		} else {
			metricsGaugeUpdate[key] = val
		}
	}
	for key, val := range metricsGaugeUpdate {
		query := "UPDATE metrics SET value = $1 WHERE id = $2"
		if err := st.exec(context.Background(), tx, query, val, key); err != nil {
			errr := tx.Rollback()
			if errr != nil {
				return fmt.Errorf("error while trying to save batch gauge: err%w;And can not rollback! err:%w",
					err, errr)
			}
			return fmt.Errorf("error while trying to save all gauge metric: %w", err)
		}
	}
	counter := 1
	valuesForInsert := make([]any, 0)
	qeryTest := `INSERT INTO metrics (id, mtype, value) VALUES `
	for key, val := range metricsGaugeInsert {
		if counter == 1 {
			qeryTest += fmt.Sprintf(" ($%d, $%d, $%d)", counter, counter+1, counter+2)
		} else {
			qeryTest += fmt.Sprintf(", ($%d, $%d, $%d)", counter, counter+1, counter+2)
		}
		valuesForInsert = append(valuesForInsert, key, typegauge, val)
		counter += 3
	}
	// qeryTest += " ON DUPLICATE KEY DO UPDATE value = EXCLUDED.value"
	if err := st.exec(context.Background(), tx, qeryTest, valuesForInsert...); err != nil {
		errr := tx.Rollback()
		if errr != nil {
			return fmt.Errorf("error while trying to save batch gauge: err%w;And can not rollback! err:%w",
				err, errr)
		}
		return fmt.Errorf("error while trying to save all gauge metric: %w", err)
	}

	// counter := 1
	// valuesForInsert := make([]any, 0)
	// qeryTest := `INSERT INTO metrics (id, mtype, value) VALUES `
	// for _, metric := range metricsGaugeInsert {
	// 	if counter == 1 {
	// 		qeryTest += fmt.Sprintf(" ($%d, $%d, $%d)", counter, counter+1, counter+2)
	// 	} else {
	// 		qeryTest += fmt.Sprintf(", ($%d, $%d, $%d)", counter, counter+1, counter+2)
	// 	}
	// 	valuesForInsert = append(valuesForInsert, metric.ID, typegauge, metric.GetValueMetric())
	// 	counter += 3
	// }
	// if err := st.exec(context.Background(), tx, qeryTest, valuesForInsert...); err != nil {
	// 	errr := tx.Rollback()
	// 	if errr != nil {
	// 		return fmt.Errorf("error while trying to save batch gauge: err%w;And can not rollback! err:%w",
	// 			err, errr)
	// 	}
	// 	return fmt.Errorf("error while trying to save all gauge metric: %w", err)
	// }
	for _, metric := range metrics {
		switch metric.MType {
		case "counter":
			continue
		case "gauge":
			continue
			// query := `INSERT INTO metrics (id, mtype, value) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET value = $4`
			// // ON CONFLICT (id) DO UPDATE SET value = $4;
			// if err := st.exec(context.Background(), tx, query, metric.ID, typegauge,
			// 	common.TypeGauge(*metric.Value)); err != nil {
			// 	errr := tx.Rollback()
			// 	if errr != nil {
			// 		return fmt.Errorf("error while trying to save gauge metric %s: err%w;And can not rollback! err:%w",
			// 			metric.ID, err, errr)
			// 	}
			// 	return fmt.Errorf("error while trying to save gauge metric %s: %w", metric.ID, err)
			// }
		default:
			st.Logger.Debugf("AddMetrics. unknown type:%s;", metric.MType)
			continue
		}
	}
	comerr := tx.Commit()
	if comerr != nil {
		return fmt.Errorf("error with commit transactiom. err:%w", err)
	}
	for _, metric := range metricsCounter {
		switch metric.MType {
		case "counter":
			_ = st.AddCounter(metric.ID, common.TypeCounter(*metric.Delta))
		default:
			continue
		}
	}

	return nil
}
