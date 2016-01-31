package viperprovider

import "github.com/spf13/viper"

// V is a dummy struct used to export viper data with the config interface
// viper uses global state (whee), so we don't store any data in this package
type V struct{}

// Get looks in the global viper for the specified key
func (V) Get(key string) (string, bool) {
	return viper.GetString(key), viper.IsSet(key)
}
