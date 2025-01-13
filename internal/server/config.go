package server

type Config struct {
	ServerAdderess      string
	FileStoregePath     string
	DatabaseDsn         string
	StoreIntervalSecond int
	Restore             bool
}
