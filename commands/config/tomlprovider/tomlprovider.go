package tomlprovider

import "github.com/BurntSushi/toml"

// P contains the data for a TOML config file
type P struct {
	data configFile
}

// New returns a loaded config file at the location specified
// and error if it's unable to load it
func New(file string) (*P, error) {
	data := new(configFile)
	_, err := toml.DecodeFile(file, data)
	if err != nil {
		return nil, err
	}
	return &P{data: *data}, nil
}

// Get returns the value at specified key if it exists
func (p *P) Get(key string) (string, bool) {
	d, ok := p.data[key]
	return d, ok
}

type configFile map[string]string
