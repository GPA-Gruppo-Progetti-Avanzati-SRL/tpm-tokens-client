package tokensclient

import (
	"errors"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/har"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-tokens-client/tokensclient/model/bearer"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-tokens-client/tokensclient/model/businessview"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-tokens-client/tokensclient/model/facts"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-tokens-client/tokensclient/model/token"
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

func NewTokensApiClient(cfg *Config, opts ...restclient.Option) (*Client, error) {
	const semLogContext = "new-tokens-api-client"
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

func DeserializeTokenContextContentResponse(resp *har.Entry) (*token.TokenContext, error) {

	const semLogContext = "tokens-api-client::deserialize-token-context-response"
	if resp == nil || resp.Response == nil || resp.Response.Content == nil || resp.Response.Content.Data == nil {
		err := errors.New("cannot deserialize null response")
		log.Error().Err(err).Msg(semLogContext)
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	var resultObj *token.TokenContext
	var err error
	switch resp.Response.Status {
	case http.StatusOK:
		resultObj, err = token.DeserializeContext(resp.Response.Content.Data)
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

func DeserializeQueryBearersContentResponse(resp *har.Entry) (*bearer.BearersQueryResponse, error) {

	const semLogContext = "tokens-api-client::deserialize-query-facts-response"
	if resp == nil || resp.Response == nil || resp.Response.Content == nil || resp.Response.Content.Data == nil {
		err := errors.New("cannot deserialize null response")
		log.Error().Err(err).Msg(semLogContext)
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	var resultObj *bearer.BearersQueryResponse
	var err error
	switch resp.Response.Status {
	case http.StatusOK:
		resultObj, err = bearer.DeserializeBearersQueryResponse(resp.Response.Content.Data)
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

func DeserializeBearerContentResponse(resp *har.Entry) (*bearer.Bearer, error) {

	const semLogContext = "tokens-api-client::deserialize-bearer-response"
	if resp == nil || resp.Response == nil || resp.Response.Content == nil || resp.Response.Content.Data == nil {
		err := errors.New("cannot deserialize null response")
		log.Error().Err(err).Msg(semLogContext)
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	var resultObj *bearer.Bearer
	var err error
	switch resp.Response.Status {
	case http.StatusOK:
		resultObj, err = bearer.DeserializeBearer(resp.Response.Content.Data)
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

func DeserializeTokenResponseBody(resp *har.Entry) (*token.Token, error) {

	const semLogContext = "tokens-api-client::deserialize-token-response"
	if resp == nil || resp.Response == nil || resp.Response.Content == nil || resp.Response.Content.Data == nil {
		err := errors.New("cannot deserialize null response")
		log.Error().Err(err).Msg(semLogContext)
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	var resultObj *token.Token
	var err error
	switch resp.Response.Status {
	case http.StatusOK:
		resultObj, err = token.DeserializeToken(resp.Response.Content.Data)
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

func DeserializeTokenTimerResponseBody(resp *har.Entry) (*token.Timer, error) {

	const semLogContext = "tokens-api-client::deserialize-token-timer-response"
	if resp == nil || resp.Response == nil || resp.Response.Content == nil || resp.Response.Content.Data == nil {
		err := errors.New("cannot deserialize null response")
		log.Error().Err(err).Msg(semLogContext)
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	var resultObj *token.Timer
	var err error
	switch resp.Response.Status {
	case http.StatusOK:
		resultObj, err = token.DeserializeTimer(resp.Response.Content.Data)
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

func DeserializeTokenTimersResponseBody(resp *har.Entry) ([]token.Timer, error) {

	const semLogContext = "tokens-api-client::deserialize-token-timers-response"
	if resp == nil || resp.Response == nil || resp.Response.Content == nil || resp.Response.Content.Data == nil {
		err := errors.New("cannot deserialize null response")
		log.Error().Err(err).Msg(semLogContext)
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	var resultObj []token.Timer
	var err error
	switch resp.Response.Status {
	case http.StatusOK:
		resultObj, err = token.DeserializeTimers(resp.Response.Content.Data)
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

func DeserializeTokenViewContentResponse(resp *har.Entry) (*businessview.Token, error) {

	const semLogContext = "tokens-api-client::deserialize-token-view-response"
	if resp == nil || resp.Response == nil || resp.Response.Content == nil || resp.Response.Content.Data == nil {
		err := errors.New("cannot deserialize null response")
		log.Error().Err(err).Msg(semLogContext)
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	var resultObj *businessview.Token
	var err error
	switch resp.Response.Status {
	case http.StatusOK:
		resultObj, err = businessview.DeserializeToken(resp.Response.Content.Data)
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

func DeserializeActorViewContentResponse(resp *har.Entry) (*businessview.Actor, error) {

	const semLogContext = "tokens-api-client::deserialize-actor-view-response"
	if resp == nil || resp.Response == nil || resp.Response.Content == nil || resp.Response.Content.Data == nil {
		err := errors.New("cannot deserialize null response")
		log.Error().Err(err).Msg(semLogContext)
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	var resultObj *businessview.Actor
	var err error
	switch resp.Response.Status {
	case http.StatusOK:
		resultObj, err = businessview.DeserializeActor(resp.Response.Content.Data)
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

func DeserializeQueryFactsContentResponse(resp *har.Entry) (*facts.FactsQueryResponse, error) {

	const semLogContext = "tokens-api-client::deserialize-query-facts-response"
	if resp == nil || resp.Response == nil || resp.Response.Content == nil || resp.Response.Content.Data == nil {
		err := errors.New("cannot deserialize null response")
		log.Error().Err(err).Msg(semLogContext)
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	var resultObj *facts.FactsQueryResponse
	var err error
	switch resp.Response.Status {
	case http.StatusOK:
		resultObj, err = facts.DeserializeFactsQueryResponse(resp.Response.Content.Data)
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

func DeserializeFactContentResponse(resp *har.Entry) (*facts.Fact, error) {

	const semLogContext = "tokens-api-client::deserialize-fact-response"
	if resp == nil || resp.Response == nil || resp.Response.Content == nil || resp.Response.Content.Data == nil {
		err := errors.New("cannot deserialize null response")
		log.Error().Err(err).Msg(semLogContext)
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	var resultObj *facts.Fact
	var err error
	switch resp.Response.Status {
	case http.StatusOK:
		resultObj, err = facts.DeserializeFact(resp.Response.Content.Data)
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

	const semLogContext = "tokens-api-client::deserialize-api-response"
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
