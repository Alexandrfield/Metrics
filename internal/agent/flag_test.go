package agent

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseFile(t *testing.T) {
	os.Clearenv()
	defer os.Clearenv()

	expectedEnvServerAdderess := "localhost:8080"
	expectedEnvSignKey := "ffab"
	expectedEnvCryptKey := "test"
	os.Setenv("ADDRESS", expectedEnvServerAdderess)
	os.Setenv("KEY", expectedEnvSignKey)
	os.Setenv("RATE_LIMIT", "1")
	os.Setenv("REPORT_INTERVAL", "2")
	os.Setenv("POLL_INTERVAL", "4")
	os.Setenv("CRYPTO_KEY", expectedEnvCryptKey)
	var config Config
	err := parseFlags(&config)
	require.NoError(t, err)

	if config.ServerAdderess != expectedEnvServerAdderess {
		t.Errorf("ServerAdderess actual:%s; expected:%s", config.ServerAdderess, expectedEnvServerAdderess)
	}
	if string(config.SignKey) != "" {
		t.Errorf("SignKey actual:%s; expected:%s", string(config.SignKey), "")
	}
	if config.PollIntervalSecond != 4 {
		t.Errorf("PollIntervalSecond actual:%d; expected:%d", config.PollIntervalSecond, 4)
	}
	if config.ReportIntervalSecond != 2 {
		t.Errorf("ReportIntervalSecond actual:%d; expected:%d", config.ReportIntervalSecond, 2)
	}
	if config.RateLimit != 1 {
		t.Errorf("RateLimit actual:%d; expected:%d", config.RateLimit, 1)
	}
	if len(config.CryptoKeyOpen) != 0 {
		t.Error("expected empty key")
	}
}

func TestParseJSONFile(t *testing.T) {
	var data string = `{"address": "localhost:8080","report_interval": "1s","poll_interval": "1s",
	 "crypto_key": "/path/to/key.pem"}`
	var conf *Config = parseJson([]byte(data))
	if conf == nil {
		t.Error("config nil!")
		return
	}

	if conf.PollIntervalSecond != 1 {
		t.Errorf("StoreIntervalSecond actual:%d, expected:%d", conf.PollIntervalSecond, 1)
	}
	if conf.ServerAdderess != "localhost:8080" {
		t.Errorf("ServerAdderess actual:%s, expected:%s", conf.ServerAdderess, "localhost:8080")
	}
}

func TestParseJSONFileError(t *testing.T) {
	data := [][]byte{
		[]byte(`{"address": "localhost:8080","report_interval": "1qwerty","poll_interval": "1s",
	 "crypto_key": "/path/to/key.pem"}`),
		[]byte(`{"address": "localhost:8080","report_interval": "1s","poll_interval": "qwerty",
	 "crypto_key": "/path/to/key.pem"}`),
		[]byte{0x00, 0x01},
	}
	for _, v := range data {
		var conf *Config = parseJson(v)
		if conf != nil {
			t.Errorf("config is not nil! data:%s", string(v))
			return
		}
	}
}
