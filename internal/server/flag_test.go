package server

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseFile(t *testing.T) {
	var config Config
	err := parseFlags(&config)
	require.NoError(t, err)

	expectedEnvDatabaseDsn, ok := os.LookupEnv("DATABASE_DSN")
	if !ok {
		expectedEnvDatabaseDsn = ""
	}
	expectedEnvFileStoregePath, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if !ok {
		expectedEnvFileStoregePath = "localStorage.dat"
	}
	expectedEnvServerAdderess, ok := os.LookupEnv("ADDRESS")
	if !ok {
		expectedEnvServerAdderess = "localhost:8080"
	}
	expectedEnvSignKey, ok := os.LookupEnv("KEY")
	if !ok {
		expectedEnvSignKey = ""
	}

	if config.DatabaseDsn != expectedEnvDatabaseDsn {
		t.Errorf("DatabaseDsn actual:%s; expected:%s", config.DatabaseDsn, expectedEnvDatabaseDsn)
	}
	if config.FileStoregePath != expectedEnvFileStoregePath {
		t.Errorf("FileStoregePath actual:%s; expected:%s", config.FileStoregePath, expectedEnvFileStoregePath)
	}
	if config.ServerAdderess != expectedEnvServerAdderess {
		t.Errorf("ServerAdderess actual:%s; expected:%s", config.ServerAdderess, expectedEnvServerAdderess)
	}
	if string(config.SignKey) != expectedEnvSignKey {
		t.Errorf("SignKey actual:%s; expected:%s", string(config.SignKey), expectedEnvSignKey)
	}
	if config.StoreIntervalSecond != 300 {
		t.Errorf("StoreIntervalSecond actual:%d; expected:%d", config.StoreIntervalSecond, 300)
	}
}
