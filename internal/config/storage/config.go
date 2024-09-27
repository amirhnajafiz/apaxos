package storage

// RedisConfig has the needed data to connect
// to a redis cache cluster.
type RedisConfig struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	Password string `koanf:"password"`
	Database int    `koanf:"database"`
}

// MySQLConfig has the needed data
// to connect to a MySQL server.
type MySQLConfig struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	User     string `koanf:"user"`
	Pass     string `konaf:"password"`
	Database string `koanf:"database"`
}
