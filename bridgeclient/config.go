package bridgeclient

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
)

const (
	TokenContextIdPathPlaceHolder = "{context-id}"
)

type EndpointDefinition struct {
	Id  string `mapstructure:"id" json:"id" yaml:"id"`
	Url string `mapstructure:"url" json:"url" yaml:"url"`
}

type HostInfo struct {
	Scheme   string `mapstructure:"scheme,omitempty" json:"scheme,omitempty" yaml:"scheme,omitempty"`
	HostName string `mapstructure:"name,omitempty" json:"name,omitempty" yaml:"name,omitempty"`
	Port     int    `mapstructure:"port,omitempty" json:"port,omitempty" yaml:"port,omitempty"`
}

// Config Note: the json serialization seems not need any inline, squash of sorts...
type Config struct {
	restclient.Config `mapstructure:",squash"  yaml:",inline"`
	Host              HostInfo             `mapstructure:"host,omitempty" json:"host,omitempty" yaml:"host,omitempty"`
	Endpoints         []EndpointDefinition `mapstructure:"endpoints" json:"endpoints" yaml:"endpoints"`
}

func (c *Config) PostProcess() error {
	return nil
}
