package server

type Config struct {
	ServerAdderess      string
	FileStoregePath     string
	DatabaseDsn         string
	SignKey             []byte
	CryptoKeySec        []byte
	StoreIntervalSecond int
	Restore             bool
}
