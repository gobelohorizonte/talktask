package config

import (
	"os"
	"strconv"

	"github.com/pkg/errors"
)

//Server ...
type Server struct {
	Host             string
	Port             string
	UseSystemdSocket bool
}

// Config ...
type Config struct {
	Server  Server
	Options struct {
		ACD struct {
			PoolSize  int
			QueueSize int
		}
	}
}

// Load configs from environment
func Load() (cfg *Config, err error) {
	cfg = &Config{}

	cfg.Server.Host = os.Getenv("GOTALK_HOST")
	cfg.Server.Port = os.Getenv("GOTALK_PORT")
	cfg.Server.UseSystemdSocket = os.Getenv("GOTALK_USE_SYSTEMD_SOCKET") == "true"

	cfg.Options.ACD.PoolSize, err = strconv.Atoi(os.Getenv("GOTALK_ACD_POOL_SIZE"))
	if err != nil {
		panic(errors.Wrap(err, "got an invalid value for 'GOTALK_ACD_POOL_SIZE'"))
	}

	cfg.Options.ACD.QueueSize, err = strconv.Atoi(os.Getenv("GOTALK_ACD_QUEUEL_SIZE"))
	if err != nil {
		panic(errors.Wrap(err, "got an invalid value for 'GOTALK_ACD_QUEUEL_SIZE'"))
	}

	return
}
