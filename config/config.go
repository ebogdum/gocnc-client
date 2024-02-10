package config

import (
	"github.com/BurntSushi/toml"
	"gocnc/utils"
	"os"
)

type Config struct {
	Server struct {
		Protocol string `toml:"protocol"`
		Hostname string `toml:"hostname"`
		Port     int    `toml:"port"`
	}
	Files []struct {
		Path        string      `toml:"path"`
		RequestPath string      `toml:"req"`
		Mode        os.FileMode `toml:"mode"`
	}
	Template struct {
		Enabled bool `toml:"enabled"`
	}
}

var conf Config

func GetConfig() Config {
	_, err := toml.DecodeFile("config.toml", &conf)
	utils.Check(err)

	return conf
}
