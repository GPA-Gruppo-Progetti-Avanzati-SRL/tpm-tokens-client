package campaignclient

import (
	"errors"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/har"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
	"github.com/rs/zerolog/log"
	"net/http"
)

const (
	ApiKeyHeaderName         = "X-Api-Key"
	RequestIdHeaderName      = "Request-id"
	LraHttpContextHeaderName = "Long-Running-Action"
	ContentTypeHeaderName    = "Content-Type"

	ContentTypeApplicationJson = "application/json"
)

type Client struct {
	host   HostInfo
	client *restclient.Client
	// harEntries []*har.Entry
}

func (c *Client) Close() {
	c.client.Close()
}

func NewCampaignApiClient(cfg *Config, opts ...restclient.Option) (*Client, error) {
	const semLogContext = "new-campaign-api-client"
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
	return &Client{client: client, host: h}, nil
}

func DeserializeCampaignContentResponse(resp *har.Entry) (*Campaign, error) {

	const semLogContext = "campaign-api-client::deserialize-campaign-response"
	if resp == nil || resp.Response == nil || resp.Response.Content == nil || resp.Response.Content.Data == nil {
		err := errors.New("cannot deserialize null response")
		log.Error().Err(err).Msg(semLogContext)
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	var resultObj *Campaign
	var err error
	switch resp.Response.Status {
	case http.StatusOK:
		resultObj, err = DeserializeCampaign(resp.Response.Content.Data)
		if err != nil {
			return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
		}

	default:
		var apiResponse ApiResponse
		apiResponse, err = DeserApiResponseFromJson(resp.Response.Content.Data)
		if err != nil {
			return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
		}
		apiResponse.StatusCode = resp.Response.Status
		err = &apiResponse
		return nil, err
	}

	return resultObj, nil
}

func DeserializeApiResponse(resp *har.Entry) (*ApiResponse, error) {

	const semLogContext = "campaign-api-client::deserialize-api-response"
	if resp == nil || resp.Response == nil || resp.Response.Content == nil || resp.Response.Content.Data == nil {
		err := errors.New("cannot deserialize null response")
		log.Error().Err(err).Msg(semLogContext)
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	apiResponse, err := DeserApiResponseFromJson(resp.Response.Content.Data)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}
	apiResponse.StatusCode = resp.Response.Status
	return &apiResponse, err
}
