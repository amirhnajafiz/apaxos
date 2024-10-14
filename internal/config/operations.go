package config

import "strconv"

// GetNodes converts the nodes pair to a simple map of strings.
func (c Config) GetNodes() map[string]string {
	hashMap := make(map[string]string)

	for _, pair := range c.Nodes {
		hashMap[pair.Key] = pair.Value
	}

	return hashMap
}

// GetBalances builds a hashmap for clients and balances with inital balance value.
func (c Config) GetBalances() map[string]int64 {
	hashMap := make(map[string]int64)

	for _, pair := range c.Clients {
		val, _ := strconv.Atoi(pair.Value)
		hashMap[pair.Key] = int64(val)
	}

	return hashMap
}

// GetClients converts the clients pair to a simple map of strings.
func (c Config) GetClients() map[string]string {
	hashMap := make(map[string]string)

	for _, pair := range c.Clients {
		hashMap[pair.Key] = pair.Value
	}

	return hashMap
}
