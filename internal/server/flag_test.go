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

func TestParseJSONFile(t *testing.T) {
	var data string = `{"address":"localhost:8080","restore":true,"store_interval":"1s","store_file":"/path/to/file.db",
	"database_dsn":"","crypto_key":"/path/to/key.pem"}`
	var conf *Config = parseJson([]byte(data))
	if conf == nil {
		t.Error("config nil!")
		return
	}
	if conf.FileStoregePath != "/path/to/file.db" {
		t.Errorf("FileStoregePath actual:%s, expected:%s", conf.FileStoregePath, "/path/to/file.db")
	}
	if conf.StoreIntervalSecond != 1 {
		t.Errorf("StoreIntervalSecond actual:%d, expected:%d", conf.StoreIntervalSecond, 1)
	}
	if conf.ServerAdderess != "localhost:8080" {
		t.Errorf("ServerAdderess actual:%s, expected:%s", conf.ServerAdderess, "localhost:8080")
	}
}
