package storage

type RedisConfig struct{}

type MySQLConfig struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	User     string `koanf:"user"`
	Pass     string `konaf:"password"`
	Database string `koanf:"database"`
}
