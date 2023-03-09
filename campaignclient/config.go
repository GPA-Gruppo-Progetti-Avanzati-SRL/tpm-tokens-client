package campaignclient

import "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"

const (
	CampaignIdPathPlaceHolder = "{campaign-id}"

	CampaignBasePath = "/api/v1/campaigns"
	CampaignQuery    = CampaignBasePath
	CampaignNew      = CampaignBasePath
	CampaignGet      = CampaignBasePath + "/" + CampaignIdPathPlaceHolder
	CampaignPut      = CampaignBasePath + "/" + CampaignIdPathPlaceHolder
	CampaignDelete   = CampaignBasePath + "/" + CampaignIdPathPlaceHolder
)

type HostInfo struct {
	Scheme   string `mapstructure:"scheme,omitempty" json:"scheme,omitempty" yaml:"scheme,omitempty"`
	HostName string `mapstructure:"name,omitempty" json:"name,omitempty" yaml:"name,omitempty"`
	Port     int    `mapstructure:"port,omitempty" json:"port,omitempty" yaml:"port,omitempty"`
}

// Config Note: the json serialization seems not need any inline, squash of sorts...
type Config struct {
	restclient.Config `mapstructure:",squash"  yaml:",inline"`
	Host              HostInfo `mapstructure:"host,omitempty" json:"host,omitempty" yaml:"host,omitempty"`
}

func (c *Config) PostProcess() error {
	return nil
}
