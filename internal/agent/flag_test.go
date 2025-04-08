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
