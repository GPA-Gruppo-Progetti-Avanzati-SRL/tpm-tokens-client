package tokensclient

import (
	"encoding/json"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/har"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-tokens-client/tokensclient/model/facts"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-tokens-client/tokensclient/model/token"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
	"strings"
)

type FactApiRequest struct {
	Id         string                 `yaml:"id,omitempty" mapstructure:"id,omitempty" json:"id,omitempty"`
	CtxId      string                 `yaml:"ctx-id,omitempty" mapstructure:"ctx-id,omitempty" json:"ctx-id,omitempty"`
	TokenId    string                 `yaml:"token-id,omitempty" mapstructure:"token-id,omitempty" json:"token-id,omitempty"`
	Properties map[string]interface{} `yaml:"properties,omitempty" mapstructure:"properties,omitempty" json:"properties,omitempty"`
	TTL        int                    `yaml:"ttl,omitempty" mapstructure:"ttl,omitempty" json:"ttl,omitempty"`
}

func (c *Client) QueryFacts(reqCtx ApiRequestContext, factsClass, factsGroup string) (*facts.FactsQueryResponse, error) {
	const semLogContext = "tpm-tokens-client::query-facts"
	log.Trace().Msg(semLogContext)

	ep := c.factsApiUrl(FactsQueryGroup, factsClass, factsGroup, "", nil)

	req, err := c.client.NewRequest(http.MethodGet, ep, nil, reqCtx.getHeaders(""), nil)
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	harEntry, err := c.client.Execute(req,
		restclient.ExecutionWithOpName("client-query-facts"),
		restclient.ExecutionWithRequestId(reqCtx.RequestId),
		restclient.ExecutionWithSpan(reqCtx.Span),
		restclient.ExecutionWithHarSpan(reqCtx.HarSpan))
	// c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeQueryFactsContentResponse(harEntry)
	return resp, err
}

func (c *Client) AddFact2Group(reqCtx ApiRequestContext, factsClass, factsGroup string, fact *FactApiRequest) (*facts.Fact, error) {
	const semLogContext = "tpm-tokens-client::add-fact-2-group"
	log.Trace().Msg(semLogContext)

	ep := c.factsApiUrl(FactAdd2Group, factsClass, factsGroup, "", nil)
	ct := ContentTypeApplicationJson

	b, err := json.Marshal(fact)
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	req, err := c.client.NewRequest(http.MethodPost, ep, b, reqCtx.getHeaders(ct), nil)
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	harEntry, err := c.client.Execute(req,
		restclient.ExecutionWithOpName("add-fact-2-group"),
		restclient.ExecutionWithRequestId(reqCtx.RequestId),
		restclient.ExecutionWithSpan(reqCtx.Span),
		restclient.ExecutionWithHarSpan(reqCtx.HarSpan))
	// c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeFactContentResponse(harEntry)
	return resp, err
}

func (c *Client) factsApiUrl(apiPath string, factsClass, factGroup, factId string, qParams []har.NameValuePair) string {
	var sb = strings.Builder{}
	sb.WriteString(c.host.Scheme)
	sb.WriteString("://")
	sb.WriteString(c.host.HostName)
	sb.WriteString(":")
	sb.WriteString(fmt.Sprint(c.host.Port))

	apiPath = strings.Replace(apiPath, FactClassPathPlaceHolder, token.WellFormTokenContextId(factsClass), 1)
	apiPath = strings.Replace(apiPath, FactGroupPathPlaceHolder, token.WellFormTokenContextId(factGroup), 1)
	apiPath = strings.Replace(apiPath, FactIdPathPlaceHolder, token.WellFormTokenId(factId), 1)
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
