package config

import (
	"flag"
	"io/ioutil"
	stdlog "log"
	"os"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
	yaml "gopkg.in/yaml.v2"
)

var (
	// EnvPrefix is prefix for identity
	EnvPrefix string
)

type LogSetting struct {
	Name             string `yaml:"name"`
	Type             string `yaml:"type"`
	MinLevel         string `yaml:"min_level"`
	ConnectionString string `yaml:"connection_string"`
}

type Configuration struct {
	Env      string
	Logs     []LogSetting `yaml:"logs"`
	Database struct {
		Username string
		Password string
		Address  string
		Type     string
		DBName   string
	}
	Redis struct {
		Address  string
		Password string
		DB       int
	}
	Nats struct {
		ClusterID string `yaml:"cluster_id"`
		Username  string
		Password  string
		Address   string
	}
	Identity struct {
		AdvertiseAddr string `yaml:"advertise_addr"`
		GRPCBind      string `yaml:"grpc_bind"`
	}
}

func New(fileName string) Configuration {
	flag.Parse()
	c := Configuration{}

	rootDirPath := os.Getenv("IDENTITY_HOME")
	if rootDirPath == "" {
		//read and parse config file
		rootDirPathStr, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			stdlog.Fatalf("config: file error: %s", err.Error())
		}
		rootDirPath = rootDirPathStr
	}
	configPath := filepath.Join(rootDirPath, "configs", fileName)
	_, err := os.Stat(configPath)
	if err != nil {
		stdlog.Fatalf("config: file error: %s", err.Error())
	}

	// config exists
	file, err := ioutil.ReadFile(filepath.Clean(configPath))
	if err != nil {
		stdlog.Fatalf("config: read file error: %s", err.Error())
	}

	err = yaml.Unmarshal(file, &c)
	if err != nil {
		stdlog.Fatal("config: yaml unmarshal error:", err)
	}

	if EnvPrefix == "" {
		stdlog.Fatal("config: env prefix not set")
	}

	if err := envconfig.Process(EnvPrefix, &c); err != nil {
		stdlog.Fatal("config: env failed:", err)
	}

	return c
}
