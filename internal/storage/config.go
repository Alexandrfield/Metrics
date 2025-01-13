package storage

type Config struct {
	FileStoregePath     string
	DatabaseDsn         string
	StoreIntervalSecond int
	Restore             bool
}
