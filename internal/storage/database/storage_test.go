package databasestorage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Alexandrfield/Metrics/internal/common"
	mock "github.com/Alexandrfield/Metrics/internal/storage/database/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestExecZeroQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBConnDB := mock.NewMockdatabaseDB(ctrl)
	stor := NewMemDatabaseStorage(&common.FakeLogger{}, "test")

	ctx := context.Background()
	err := stor.exec(ctx, mockBConnDB, "")
	if err != nil {
		t.Errorf("expected nil error. actual:%s", err)
	}
}

func TestExec(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBConnDB := mock.NewMockdatabaseDB(ctrl)
	mockBConnDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
	stor := NewMemDatabaseStorage(&common.FakeLogger{}, "test")

	ctx := context.Background()
	query := `test query`
	err := stor.exec(ctx, mockBConnDB, query, "")
	if err != nil {
		t.Errorf("expected nil error. actual:%s", err)
	}
}
func TestExecErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBConnDB := mock.NewMockdatabaseDB(ctrl)
	errT := fmt.Errorf("test error")
	mockBConnDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errT).Times(1)
	stor := NewMemDatabaseStorage(&common.FakeLogger{}, "test")

	ctx := context.Background()
	query := `test query`
	err := stor.exec(ctx, mockBConnDB, query, "")
	if err == nil {
		t.Errorf("expected not nil error. actual:%s", err)
	}
	stor.Close()
}

func TestExecT(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBConnDB := mock.NewMockdatabaseDB(ctrl)
	mockBConnDB.EXPECT().QueryRowContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
	//stor := NewMemDatabaseStorage(&common.FakeLogger{}, "test")

	ctx := context.Background()
	query := `test query`
	r := mockBConnDB.QueryRowContext(ctx, query, "")
	if r != nil {
		t.Errorf("expected not nil error. ")
	}
}

func TestAddCounterMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBConnDB := mock.NewMockdatabaseDB(ctrl)
	mockBConnDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
	stor := NewMemDatabaseStorage(&common.FakeLogger{}, "test")

	countersMap := make(map[string]common.TypeCounter)
	countersMap["testCounter1"] = common.TypeCounter(4)
	countersMap["testCounter2"] = common.TypeCounter(67)
	err := stor.addCounterMetrics(mockBConnDB, countersMap)
	stor.Close()
	require.NoError(t, err)
}

func TestAddGaugeMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBConnDB := mock.NewMockdatabaseDB(ctrl)
	mockBConnDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
	stor := NewMemDatabaseStorage(&common.FakeLogger{}, "test")

	countersMap := make(map[string]common.TypeGauge)
	countersMap["testCounter1"] = common.TypeGauge(4.4)
	countersMap["testCounter2"] = common.TypeGauge(56.3)
	err := stor.addGaugeMetrics(mockBConnDB, countersMap)
	require.NoError(t, err)
}

func TestCreateTables(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbsql := mock.NewMockdbSql(ctrl)
	dbsql.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
	stor := NewMemDatabaseStorage(&common.FakeLogger{}, "test")
	stor.db = dbsql
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := stor.createTable(ctx)
	require.NoError(t, err)
}

func TestCreateTablesErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbsql := mock.NewMockdbSql(ctrl)
	errT := fmt.Errorf("test error")
	dbsql.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errT)
	stor := NewMemDatabaseStorage(&common.FakeLogger{}, "test")
	stor.db = dbsql
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := stor.createTable(ctx)
	if err == nil {
		t.Errorf("create database. expected err:%s; actual:nil", errT)
	}
}

func TestPingDatabase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbsql := mock.NewMockdbSql(ctrl)
	dbsql.EXPECT().PingContext(gomock.Any()).Return(nil)
	stor := NewMemDatabaseStorage(&common.FakeLogger{}, "test")
	stor.db = dbsql
	st := stor.PingDatabase()
	if st == false {
		t.Errorf("Ping database. expected check: true; actual:%t", st)
	}
}

func TestPingDatabaseErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbsql := mock.NewMockdbSql(ctrl)
	errT := fmt.Errorf("test error")
	dbsql.EXPECT().PingContext(gomock.Any()).Return(errT)
	stor := NewMemDatabaseStorage(&common.FakeLogger{}, "test")
	stor.db = dbsql
	st := stor.PingDatabase()
	if st == true {
		t.Errorf("Ping database. expected check: false; actual:%t", st)
	}
}

func TestClose(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbsql := mock.NewMockdbSql(ctrl)
	dbsql.EXPECT().Close().Return(nil)
	stor := NewMemDatabaseStorage(&common.FakeLogger{}, "test")
	stor.db = dbsql
	stor.Close()
}
func TestExp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbsql := mock.NewMockdbSql(ctrl)
	dbsql.EXPECT().Close().Return(nil)
	dbsql.EXPECT().Begin().Return(nil, nil)
	dbsql.EXPECT().QueryContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
	dbsql.EXPECT().QueryRowContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	stor := NewMemDatabaseStorage(&common.FakeLogger{}, "test")
	stor.db = dbsql
	stor.db.Begin()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, _ = stor.db.QueryContext(ctx, "test", "er")
	_ = stor.db.QueryRowContext(ctx, "test", "er")
	stor.Close()
}

func TestStart(t *testing.T) {
	stor := NewMemDatabaseStorage(&common.FakeLogger{}, "test")
	err := stor.Start()
	if err == nil {
		t.Error("Start. Expected not nil err")
	}
}
