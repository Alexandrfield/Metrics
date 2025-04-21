package agent

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseFile(t *testing.T) {
	var config Config
	err := parseFlags(&config)
	require.NoError(t, err)

	expectedEnvServerAdderess, ok := os.LookupEnv("ADDRESS")
	if !ok {
		expectedEnvServerAdderess = "localhost:8080"
	}
	expectedEnvSignKey, ok := os.LookupEnv("KEY")
	if !ok {
		expectedEnvSignKey = ""
	}

	if config.ServerAdderess != expectedEnvServerAdderess {
		t.Errorf("ServerAdderess actual:%s; expected:%s", config.ServerAdderess, expectedEnvServerAdderess)
	}
	if string(config.SignKey) != expectedEnvSignKey {
		t.Errorf("SignKey actual:%s; expected:%s", string(config.SignKey), expectedEnvSignKey)
	}
	expectedPollIntervalSecond := 2
	if config.PollIntervalSecond != expectedPollIntervalSecond {
		t.Errorf("StoreIntervalSecond actual:%d; expected:%d", config.PollIntervalSecond, expectedPollIntervalSecond)
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
