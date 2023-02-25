package bridgeclient

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
)

type LinkedService struct {
	cfg *Config
}

func NewInstanceWithConfig(cfg *Config) (*LinkedService, error) {
	lks := &LinkedService{cfg: cfg}
	return lks, nil
}

func (lks *LinkedService) NewClient(opts ...restclient.Option) (*Client, error) {

	client, err := NewClient(lks.cfg, opts...)
	return client, err
}
