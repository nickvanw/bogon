package config

import "os"

// providers contians all of the providers that we're going to look at getting
// config data. By default, just looks in the environment
// global state rules
var providers = []configProvider{envProvider}

// RegisterProvider adds a new config provider at the end of the priority
func RegisterProvider(c configProvider) {
	providers = append(providers, c)
}

// Get returns the value at the specified string, and true if it exists
func Get(key string) (string, bool) {
	for _, v := range providers {
		if ret, ok := v.Get(key); ok {
			return ret, ok
		}
	}
	return "", false
}

type configProvider interface {
	Get(string) (string, bool)
}

var envProvider = EnvProvider{}

// EnvProvider searches the environment for the key
type EnvProvider struct{}

// Get returns the environment variable at string, and true if it exists
func (e EnvProvider) Get(key string) (string, bool) {
	return os.LookupEnv(key)
}
