package config

import (
	"io"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

// define the struct for config
type AppOptions struct {
	Server *ServerConfig `yaml:"server"`
	Log    *LogConfig    `yaml:"log"`

	RedisConfig *redis.Options `yaml:"redis_config"`
}

func (opt *AppOptions) FillWithDefaults() {
	opt.Server.FillWithDefaults()
}

type ServerConfig struct {
	Addr string `yaml:"addr"` // bind address include port

	RTimeout       int `yaml:"request_timeout"`  // request timeout, in second
	WTimeout       int `yaml:"response_timeout"` // response timeout, in second
	DTimeout       int `yaml:"idle_timeout"`     // http connection idle timeout, in second
	MaxHeaderBytes int `yaml:"max_header_bytes"` // unit in byte
}

func (cfg *ServerConfig) FillWithDefaults() {
	if cfg == nil {
		return
	}
	if cfg.RTimeout == 0 {
		cfg.RTimeout = 3
	}
	if cfg.WTimeout == 0 {
		cfg.WTimeout = 10
	}
	if cfg.DTimeout == 0 {
		cfg.DTimeout = 30
	}
	if cfg.MaxHeaderBytes == 0 {
		cfg.MaxHeaderBytes = 1 << 20
	}
}

type LogConfig struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"` // valid values [stdout|stderr|path/to/file]
}

func (cfg *LogConfig) GetLevel() (logrus.Level, error) {
	return logrus.ParseLevel(cfg.Level)
}

func (cfg *LogConfig) GetWriter() (io.WriteCloser, error) {
	if len(cfg.Output) == 0 {
		return os.Stderr, nil
	}
	switch cfg.Output {
	case "stderr":
		return os.Stderr, nil
	case "stdout":
		return os.Stdout, nil
	default:
		// no need to close currently
		return os.OpenFile(cfg.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	}
}
