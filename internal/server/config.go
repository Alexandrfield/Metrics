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

type configJSON struct {
	ServerAdderess      string `json:"address"`
	FileStoregePath     string `json:"store_file"`
	DatabaseDsn         string `json:"database_dsn"`
	SignKey             string `json:"sign_key"`
	CryptoKeySec        string `json:"crypto_key"`
	StoreIntervalSecond string `json:"store_interval"`
	Restore             bool   `json:"restore"`
}
