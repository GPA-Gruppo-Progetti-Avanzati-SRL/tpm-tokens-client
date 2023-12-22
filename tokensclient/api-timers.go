package tokensclient

import (
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/har"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-tokens-client/tokensclient/model/token"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) CreateTimers(reqCtx ApiRequestContext, ctxId string, tokId string) ([]token.Timer, error) {
	const semLogContext = "tpm-tokens-client::post-create-timers"

	ep := c.timerApiUrl(TokenTimerCreate, ctxId, tokId, nil)

	req, err := c.client.NewRequest(http.MethodPost, ep, nil, reqCtx.getHeaders(""), nil)
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	harEntry, err := c.client.Execute(req,
		restclient.ExecutionWithOpName("client-token-create-timer"),
		restclient.ExecutionWithRequestId(reqCtx.RequestId),
		// restclient.ExecutionWithLraId(reqCtx.LRAId),
		restclient.ExecutionWithSpan(reqCtx.Span),
		restclient.ExecutionWithHarSpan(reqCtx.HarSpan))
	// c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeTokenTimersResponseBody(harEntry)
	return resp, err
}

func (c *Client) DeleteTimers(reqCtx ApiRequestContext, ctxId string, tokId string) (*ApiResponse, error) {
	const semLogContext = "tpm-tokens-client::post-delete-timers"

	ep := c.timerApiUrl(TokenTimersDelete, ctxId, tokId, nil)

	req, err := c.client.NewRequest(http.MethodDelete, ep, nil, reqCtx.getHeaders(""), nil)
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	harEntry, err := c.client.Execute(req,
		restclient.ExecutionWithOpName("client-token-delete-timers"),
		restclient.ExecutionWithRequestId(reqCtx.RequestId),
		// restclient.ExecutionWithLraId(reqCtx.LRAId),
		restclient.ExecutionWithSpan(reqCtx.Span),
		restclient.ExecutionWithHarSpan(reqCtx.HarSpan))
	// c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeApiResponse(harEntry)
	return resp, err
}

func (c *Client) timerApiUrl(apiPath string, ctxId string, tokenId string, qParams []har.NameValuePair) string {
	var sb = strings.Builder{}
	sb.WriteString(c.host.Scheme)
	sb.WriteString("://")
	sb.WriteString(c.host.HostName)
	sb.WriteString(":")
	sb.WriteString(fmt.Sprint(c.host.Port))

	apiPath = strings.Replace(apiPath, TokenContextIdPathPlaceHolder, token.WellFormTokenContextId(ctxId), 1)
	apiPath = strings.Replace(apiPath, TokenIdPathPlaceHolder, token.WellFormTokenId(tokenId), 1)
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
