package storage

// Redis config contains redis host, port, and password.
// Although the default database is 0, we can set it to other dbs.
type Redis struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	Password string `koanf:"password"`
	Database int    `koanf:"database"`
}

// MongoDB config contains a uri string
// which can be provided by a MongoDB cluster.
// It also includes a database parameter.
type MongoDB struct {
	URI      string `koanf:"uri"`
	Database string `koanf:"database"`
}
