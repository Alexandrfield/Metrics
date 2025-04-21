package agent

type configJSON struct {
	ServerAdderess       string `json:"address"`
	SignKey              string `json:"sign_key"`
	CryptoKeyOpen        string `json:"crypto_key"`
	PollIntervalSecond   string `json:"poll_interval"`
	ReportIntervalSecond string `json:"report_interval"`
	RateLimit            int    `json:"rate_limit"`
}

type Config struct {
	ServerAdderess       string
	SignKey              []byte
	CryptoKeyOpen        []byte
	PollIntervalSecond   int
	ReportIntervalSecond int
	RateLimit            int
}
