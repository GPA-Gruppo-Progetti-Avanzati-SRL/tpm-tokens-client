package bridgeclient

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
	"github.com/rs/zerolog/log"
	"strings"
)

const (
	ApiKeyHeaderName      = "X-Api-Key"
	RequestIdHeaderName   = "Request-id"
	ContentTypeHeaderName = "Content-Type"

	ContentTypeApplicationJson = "application/json"
)

type Client struct {
	host      HostInfo
	endpoints []EndpointDefinition
	client    *restclient.Client
}

func (c *Client) Close() {
	c.client.Close()
}

func NewClient(cfg *Config, opts ...restclient.Option) (*Client, error) {
	const semLogContext = "bridge-client::new"
	client := restclient.NewClient(&cfg.Config, opts...)

	h := cfg.Host
	if h.Scheme == "" {
		h.Scheme = "http"
	}

	if h.HostName == "" {
		h.HostName = "localhost"
	}

	if h.Port == 0 {
		switch h.Scheme {
		case "http":
			h.Port = 80
		case "https":
			h.Port = 443
		default:
			log.Error().Str("scheme", h.Scheme).Msg(semLogContext + " invalid scheme...reverting to http...")
		}
	}

	log.Trace().Str("scheme", h.Scheme).Int("port", h.Port).Str("host-name", h.HostName).Msg(semLogContext)
	return &Client{client: client, host: h, endpoints: cfg.Endpoints}, nil
}

func (cli *Client) findEndpointPathById(endpointId string) string {

	endpointId = strings.ToLower(endpointId)
	for _, ep := range cli.endpoints {
		if ep.Id == endpointId {
			return ep.Url
		}
	}

	return ""
}
