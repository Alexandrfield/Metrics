package server

type Config struct {
	ServerAdderess      string
	FileStoregePath     string
	DatabaseDsn         string
	SignKey             []byte
	StoreIntervalSecond int
	Restore             bool
}
