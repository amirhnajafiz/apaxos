package mongodb

// Config contains a uri string
// which can be provided by a MongoDB cluster.
// It also includes a database parameter.
type Config struct {
	URI      string `koanf:"uri"`
	Database string `koanf:"database"`
}
