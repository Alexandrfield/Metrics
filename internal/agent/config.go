package agent

type Config struct {
	ServerAdderess       string
	SignKey              []byte
	PollIntervalSecond   int
	ReportIntervalSecond int
	RateLimit            int
}
