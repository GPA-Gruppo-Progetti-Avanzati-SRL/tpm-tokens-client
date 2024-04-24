package tokensclient

import (
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

func (c *Client) QueryFacts(reqCtx ApiRequestContext, factsGroup string) (*facts.FactsQueryResponse, error) {
	const semLogContext = "tpm-tokens-client::query-facts"
	log.Trace().Msg(semLogContext)

	ep := c.factsApiUrl(FactsQueryGroup, factsGroup, "", nil)

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

func (c *Client) factsApiUrl(apiPath string, factGroup, factId string, qParams []har.NameValuePair) string {
	var sb = strings.Builder{}
	sb.WriteString(c.host.Scheme)
	sb.WriteString("://")
	sb.WriteString(c.host.HostName)
	sb.WriteString(":")
	sb.WriteString(fmt.Sprint(c.host.Port))

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
