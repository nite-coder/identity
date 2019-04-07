package config

import (
	"flag"
	"io/ioutil"
	stdlog "log"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
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

func New(fileName string) *Configuration {
	flag.Parse()
	c := Configuration{}

	//read and parse config file
	rootDirPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		stdlog.Fatalf("config: file error: %s", err.Error())
	}
	configPath := filepath.Join(rootDirPath, fileName)
	_, err = os.Stat(configPath)
	if err != nil {
		stdlog.Fatalf("config: file error: %s", err.Error())
	}

	// config exists
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		stdlog.Fatalf("config: read file error: %s", err.Error())
	}

	err = yaml.Unmarshal(file, &c)
	if err != nil {
		stdlog.Fatal("config: yaml unmarshal error:", err)
	}

	return &c
}
