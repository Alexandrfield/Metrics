package agent

type Config struct {
	ServerAdderess       string
	SignKey              []byte
	CryptoKeyOpen        []byte
	PollIntervalSecond   int
	ReportIntervalSecond int
	RateLimit            int
}
