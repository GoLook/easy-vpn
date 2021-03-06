package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Provider         string              `toml:"provider"`
	PrivateKeyFile   string              `toml:"ssh_private_key"`
	PublicKeyFile    string              `toml:"ssh_public_key"`
	SelfDestructFile string              `toml:"self_destruct"`
	Sleep            int                 `toml:"sleeptime"`
	Providers        map[string]Provider `toml:"providers"`
	Options          Options             `toml:"options"`
}

type Provider struct {
	ApiKey string `toml:"api_key"`
	Region string `toml:"region"`
	Size   string `toml:"size"`
	OS     string `toml:"os"`
}

type Options struct {
	Uptime      int        `toml:"max_uptime"`
	Autoconnect bool       `toml:"vpn_autoconnect"`
	ConnectCmd  [][]string `toml:"autoconnect_cmd"`
}

func LoadConfiguration(filename string) (config *Config, err error) {
	if _, err = toml.DecodeFile(filename, &config); err != nil {
		return nil, err
	}
	return
}
