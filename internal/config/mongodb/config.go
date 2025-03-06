package mongodb

// Config contains a uri string which can be provided by a MongoDB cluster.
// It also includes a database name option to chose MongoDB database.
type Config struct {
	URI      string `koanf:"uri"`
	Database string `koanf:"database"`
}
