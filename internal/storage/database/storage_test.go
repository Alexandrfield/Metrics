package databasestorage

import (
	"context"
	"fmt"
	"testing"

	"github.com/Alexandrfield/Metrics/internal/common"
	mock "github.com/Alexandrfield/Metrics/internal/storage/database/mock"
	"github.com/golang/mock/gomock"
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
