package server

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseFile(t *testing.T) {
	os.Clearenv()
	defer os.Clearenv()

	expectedEnvServerAdderess := "localhost:8080"
	expectedEnvSignKey := "testKey"
	expectedEnvCryptKey := "test"
	expectedEnvDatabaseDsn := "qwerty"
	os.Setenv("ADDRESS", expectedEnvServerAdderess)
	os.Setenv("DATABASE_DSN", expectedEnvDatabaseDsn)
	os.Setenv("FILE_STORAGE_PATH", "test")
	os.Setenv("KEY", expectedEnvSignKey)
	os.Setenv("RATE_LIMIT", "1")
	os.Setenv("REPORT_INTERVAL", "2")
	os.Setenv("POLL_INTERVAL", "4")
	os.Setenv("CRYPTO_KEY", expectedEnvCryptKey)
	var config Config
	err := parseFlags(&config)
	require.NoError(t, err)

	if config.DatabaseDsn != expectedEnvDatabaseDsn {
		t.Errorf("DatabaseDsn actual:%s; expected:%s", config.DatabaseDsn, expectedEnvDatabaseDsn)
	}
	if config.FileStoregePath != "" {
		t.Errorf("FileStoregePath actual:%s; expected:%s", config.FileStoregePath, "")
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
	if len(config.CryptoKeySec) != 0 {
		t.Error("expected empty key")
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
