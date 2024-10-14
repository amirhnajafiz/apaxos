package config

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/f24-cse535/apaxos/internal/config/grpc"
	"github.com/f24-cse535/apaxos/internal/config/mongodb"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/tidwall/pretty"
)

// Prefix indicates environment variables prefix.
const Prefix = "apax_"

// Config struct is a module that stores system configs.
// For each node, we have a unique node_id and a client to manage.
// Other configs include gRPC, MongoDB, and other nodes gRPC addresses.
type Config struct {
	NodeID   string `koanf:"node_id"`  // a unique id for each node
	Client   string `koanf:"client"`   // the client id for each node
	Majority int    `koanf:"majority"` // number of nodes to consider as majority

	WorkersInterval int    `koanf:"workers_interval"` // node jobs' interval in seconds
	LogLevel        string `koanf:"log_level"`        // node logging level (debug, info, warn, error, panic, fatal)

	Nodes   []Pair `koanf:"nodes"`   // a map of all nodes and addresses
	Clients []Pair `koanf:"clients"` // a map of all clients and nodes

	GRPC    grpc.Config    `koanf:"grpc"`    // gRPC configs
	MongoDB mongodb.Config `koanf:"mongodb"` // MongoDB configs
}

// New reads configuration with koanf, by loading a yaml config path into the Config struct.
func New(path string) Config {
	var instance Config

	k := koanf.New(".")

	// load default configuration from file
	if err := k.Load(structs.Provider(Default(), "koanf"), nil); err != nil {
		log.Fatalf("error loading default: %s", err)
	}

	// load configuration from file
	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		log.Printf("error loading config.yml: %s", err)
	}

	// load environment variables
	if err := k.Load(env.Provider(Prefix, ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, Prefix)), "__", ".")
	}), nil); err != nil {
		log.Printf("error loading environment variables: %s", err)
	}

	if err := k.Unmarshal("", &instance); err != nil {
		log.Fatalf("error unmarshalling config: %s", err)
	}

	indent, err := json.MarshalIndent(instance, "", "\t")
	if err != nil {
		log.Fatalf("error marshaling config to json: %s", err)
	}

	indent = pretty.Color(indent, nil)
	tmpl := `
	================ Loaded Configuration ================
	%s
	=============================================
	`
	log.Printf(tmpl, string(indent))

	return instance
}
