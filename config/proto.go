package config

import (
	"io"
	"os"
	"runtime"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

// define the struct for config
type AppOptions struct {
	Server *ServerConfig `yaml:"server"`
	Log    LogConfig     `yaml:"log"`

	Redis *RedisConfig `yaml:"redis"`

	FaceAI FaceAIConfig `yaml:"face_ai"`

	OSS OSSConfig `yaml:"oss"`

	STS STSConfig `yaml:"sts"`
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
		return os.Stdout, nil
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

// A config of go redis
type RedisConfig struct {
	Network              string `yaml:"network"`
	Addr                 string `yaml:"addr"`
	Passwd               string `yaml:"password"`
	DB                   int    `yaml:"database"`
	DialTimeout          int    `yaml:"dial_timeout"`
	ReadTimeout          int    `yaml:"read_timeout"`
	WriteTimeout         int    `yaml:"write_timeout"`
	PoolSize             int    `yaml:"pool_size"`
	PoolTimeout          int    `yaml:"pool_timeout"`
	MinIdleConns         int    `yaml:"min_idle_conns"`
	MaxRetries           int    `yaml:"max_retries"`
	TraceIncludeNotFound bool   `yaml:"trace_include_not_found"`
}

func (cfg *RedisConfig) FillWithDefaults() {
	maxCPU := runtime.NumCPU()

	if cfg.DialTimeout <= 0 || cfg.DialTimeout > 1000*maxCPU {
		cfg.DialTimeout = 1000
	}

	if cfg.ReadTimeout <= 0 || cfg.ReadTimeout > 1000*maxCPU {
		cfg.ReadTimeout = 1000
	}

	if cfg.WriteTimeout <= 0 || cfg.WriteTimeout > 3000*maxCPU {
		cfg.WriteTimeout = 3000
	}

	if cfg.PoolSize <= 0 {
		cfg.PoolSize = 10 * maxCPU
	}

	if cfg.PoolTimeout <= 0 || cfg.PoolTimeout > 2*maxCPU {
		cfg.PoolTimeout = 2
	}

	if cfg.MinIdleConns <= 0 || cfg.MinIdleConns > 3*maxCPU {
		cfg.MinIdleConns = 3
	}

	if cfg.MaxRetries < 0 || cfg.MaxRetries > 1*maxCPU {
		cfg.MaxRetries = 1
	}
}

func (cfg *RedisConfig) ToOptions() *redis.Options {
	return &redis.Options{
		Network:      cfg.Network,
		Addr:         cfg.Addr,
		Password:     cfg.Passwd,
		DB:           cfg.DB,
		DialTimeout:  time.Duration(cfg.DialTimeout) * time.Millisecond,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Millisecond,
		PoolSize:     cfg.PoolSize,
		PoolTimeout:  time.Duration(cfg.PoolTimeout) * time.Second,
		MinIdleConns: cfg.MinIdleConns,
		MaxRetries:   cfg.MaxRetries,
	}
}

type FaceAIConfig struct {
	Ak       string `yaml:"ak"`
	Sk       string `yaml:"sk"`
	Endpoint string `yaml:"endpoint"`
}

type OSSConfig struct {
	PublicEndpoint       string `yaml:"public_endpoint"`
	InternalEndpoint     string `yaml:"internal_endpoint"`
	PublicBucketEndpoint string `yaml:"public_bucket_endpoint"`
	BucketName           string `yaml:"bucket_name"`
	Ak                   string `yaml:"ak"`
	Sk                   string `yaml:"sk"`
	PathPrefix           string `yaml:"path_prefix"`
}

type STSConfig struct {
	RamAK string `yaml:"ram_ak"`
	RamSK string `yaml:"ram_sk"`
	ARN   string `yaml:"arn"`
}
