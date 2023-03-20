package tokensclient

import (
	"encoding/json"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/har"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
	"strings"
)

type BearerApiRequest struct {
	Origin     string                 `yaml:"origin,omitempty" mapstructure:"origin,omitempty" json:"origin,omitempty"`
	TokenRefs  []TokenRef             `yaml:"tok-refs,omitempty" mapstructure:"tok-refs,omitempty" json:"tok-refs,omitempty"`
	Properties map[string]interface{} `yaml:"properties,omitempty" mapstructure:"properties,omitempty" json:"properties,omitempty"`
	TTL        int                    `yaml:"ttl,omitempty" mapstructure:"ttl,omitempty" json:"ttl,omitempty"`
}

func (c *Client) GetBearerInContext(reqCtx ApiRequestContext, bearerId, ctxId string) (*Bearer, error) {
	const semLogContext = "tpm-tokens-client::get-bearer-in-ctx"
	log.Trace().Msg(semLogContext)

	ep := c.bearerApiUrl(BearerContextGet, bearerId, ctxId, "", nil)

	req, err := c.client.NewRequest(http.MethodGet, ep, nil, reqCtx.getHeaders(""), nil)
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	harEntry, err := c.client.Execute(req,
		restclient.ExecutionWithOpName("client-get-bearer-in-ctx"),
		restclient.ExecutionWithRequestId(reqCtx.RequestId),
		restclient.ExecutionWithSpan(reqCtx.Span),
		restclient.ExecutionWithHarSpan(reqCtx.HarSpan))
	// c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeBearerContentResponse(harEntry)
	return resp, err
}

func (c *Client) AddBearer2Context(reqCtx ApiRequestContext, bearerId, ctxId string, bearer *BearerApiRequest, ct string) (*Bearer, error) {
	const semLogContext = "tpm-tokens-client::add-bearer-2-ctx"
	log.Trace().Msg(semLogContext)

	ep := c.bearerApiUrl(BearerContextPost, bearerId, ctxId, "", nil)

	if ct == "" {
		ct = ContentTypeApplicationJson
	}

	switch ct {
	case ContentTypeApplicationJson:
	default:
		log.Warn().Str("content-type", ct).Msg(semLogContext + " unsupported content-type...using json")
		ct = ContentTypeApplicationJson
	}

	b, err := json.Marshal(bearer)
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	req, err := c.client.NewRequest(http.MethodPost, ep, b, reqCtx.getHeaders(ct), nil)
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	harEntry, err := c.client.Execute(req,
		restclient.ExecutionWithOpName("add-bearer-2-ctx"),
		restclient.ExecutionWithRequestId(reqCtx.RequestId),
		restclient.ExecutionWithSpan(reqCtx.Span),
		restclient.ExecutionWithHarSpan(reqCtx.HarSpan))
	// c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeBearerContentResponse(harEntry)
	return resp, err
}

func (c *Client) UpdateBearerInContext(reqCtx ApiRequestContext, bearerId, ctxId string, bearer *BearerApiRequest, ct string) (*Bearer, error) {
	const semLogContext = "tpm-tokens-client::update-bearer-in-ctx"
	log.Trace().Msg(semLogContext)

	ep := c.bearerApiUrl(BearerContextPut, bearerId, ctxId, "", nil)

	if ct == "" {
		ct = ContentTypeApplicationJson
	}

	switch ct {
	case ContentTypeApplicationJson:
	default:
		log.Warn().Str("content-type", ct).Msg(semLogContext + " unsupported content-type...using json")
		ct = ContentTypeApplicationJson
	}

	b, err := json.Marshal(bearer)
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	req, err := c.client.NewRequest(http.MethodPut, ep, b, reqCtx.getHeaders(ct), nil)
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	harEntry, err := c.client.Execute(req,
		restclient.ExecutionWithOpName("update-bearer-in-ctx"),
		restclient.ExecutionWithRequestId(reqCtx.RequestId),
		restclient.ExecutionWithSpan(reqCtx.Span),
		restclient.ExecutionWithHarSpan(reqCtx.HarSpan))
	// c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeBearerContentResponse(harEntry)
	return resp, err
}

func (c *Client) RemoveBearerFromContext(reqCtx ApiRequestContext, bearerId, ctxId string) (*Bearer, error) {
	const semLogContext = "tpm-tokens-client::remove-bearer-from-ctx"
	log.Trace().Msg(semLogContext)

	ep := c.bearerApiUrl(BearerContextDelete, bearerId, ctxId, "", nil)

	req, err := c.client.NewRequest(http.MethodDelete, ep, nil, reqCtx.getHeaders(""), nil)
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	harEntry, err := c.client.Execute(req,
		restclient.ExecutionWithOpName("remove-bearer-from-ctx"),
		restclient.ExecutionWithRequestId(reqCtx.RequestId),
		restclient.ExecutionWithSpan(reqCtx.Span),
		restclient.ExecutionWithHarSpan(reqCtx.HarSpan))
	// c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeBearerContentResponse(harEntry)
	return resp, err
}

func (c *Client) AddToken2BearerInContext(reqCtx ApiRequestContext, bearerId, ctxId, tokId string, role string) (*Bearer, error) {
	const semLogContext = "tpm-tokens-client::add-token-2-bearer-in-ctx"
	log.Trace().Msg(semLogContext)

	ep := c.bearerApiUrl(AddToken2BearerInContextPost, bearerId, ctxId, tokId, []har.NameValuePair{{Name: "role", Value: role}})

	req, err := c.client.NewRequest(http.MethodPost, ep, nil, reqCtx.getHeaders(""), nil)
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	harEntry, err := c.client.Execute(req,
		restclient.ExecutionWithOpName("add-token-2-bearer-in-ctx"),
		restclient.ExecutionWithRequestId(reqCtx.RequestId),
		restclient.ExecutionWithSpan(reqCtx.Span),
		restclient.ExecutionWithHarSpan(reqCtx.HarSpan))
	// c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeBearerContentResponse(harEntry)
	return resp, err
}

func (c *Client) RemoveTokenFromBearerInContext(reqCtx ApiRequestContext, bearerId, ctxId, tokId string) (*Bearer, error) {
	const semLogContext = "tpm-tokens-client::remove-token-from-bearer-in-ctx"
	log.Trace().Msg(semLogContext)

	ep := c.bearerApiUrl(RemoveTokenFromBearerInContextDelete, bearerId, ctxId, tokId, nil)

	req, err := c.client.NewRequest(http.MethodDelete, ep, nil, reqCtx.getHeaders(""), nil)
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	harEntry, err := c.client.Execute(req,
		restclient.ExecutionWithOpName("remove-token-from-bearer-in-ctx"),
		restclient.ExecutionWithRequestId(reqCtx.RequestId),
		restclient.ExecutionWithSpan(reqCtx.Span),
		restclient.ExecutionWithHarSpan(reqCtx.HarSpan))
	// c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeBearerContentResponse(harEntry)
	return resp, err
}

func (c *Client) bearerApiUrl(apiPath string, bearerId, ctxId, tokId string, qParams []har.NameValuePair) string {
	var sb = strings.Builder{}
	sb.WriteString(c.host.Scheme)
	sb.WriteString("://")
	sb.WriteString(c.host.HostName)
	sb.WriteString(":")
	sb.WriteString(fmt.Sprint(c.host.Port))

	apiPath = strings.Replace(apiPath, TokenContextIdPathPlaceHolder, ctxId, 1)
	apiPath = strings.Replace(apiPath, BearerIdPathPlaceHolder, bearerId, 1)
	apiPath = strings.Replace(apiPath, TokenIdPathPlaceHolder, tokId, 1)
	sb.WriteString(apiPath)

	if len(qParams) > 0 {
		sb.WriteString("?")
		for i, qp := range qParams {
			if i > 0 {
				sb.WriteString("&")
			}
			sb.WriteString(qp.Name)
			sb.WriteString("=")
			sb.WriteString(url.QueryEscape(qp.Value))
		}
	}
	return sb.String()
}
