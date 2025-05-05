package server

import "net"

type Config struct {
	ServerAdderess      string
	FileStoregePath     string
	DatabaseDsn         string
	trustedSubnet       *net.IPNet
	SignKey             []byte
	CryptoKeySec        []byte
	StoreIntervalSecond int
	Restore             bool
}

func (conf *Config) IsAllowedAddres(addr net.IP) bool {
	if conf.trustedSubnet == nil {
		return true
	}
	return conf.trustedSubnet.Contains(addr)
}

type configJSON struct {
	ServerAdderess      string `json:"address"`
	FileStoregePath     string `json:"store_file"`
	DatabaseDsn         string `json:"database_dsn"`
	SignKey             string `json:"sign_key"`
	CryptoKeySec        string `json:"crypto_key"`
	StoreIntervalSecond string `json:"store_interval"`
	TrustedSubnet       string `json:"trusted_subnet"`
	Restore             bool   `json:"restore"`
}
