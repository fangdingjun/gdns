package main

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

type conf struct {
	UpstreamServers  []string `yaml:"upstream_servers"`
	BootstrapServers []string `yaml:"bootstrap_servers"`
	Listen           []listen `yaml:"listen"`
	UpstreamTimeout  int      `yaml:"upstream_timeout"`
	UpstreamInsecure bool     `yaml:"upstream_insecure"`
}

type listen struct {
	Addr string `yaml:"addr"`
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
}

func loadConfig(f string) (*conf, error) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	var c conf
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
