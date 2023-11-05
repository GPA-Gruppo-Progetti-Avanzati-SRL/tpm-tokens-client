package tokensclient

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-tokens-client/tokensclient/model/token"
	"net/http"
)

func (c *Client) CreateTimer(reqCtx ApiRequestContext, ctxId string, tokId string) (*token.Timer, error) {
	const semLogContext = "tpm-tokens-client::post-create-timer"

	ep := c.tokenApiUrl(TokenTimerCreate, ctxId, tokId, nil)

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

	resp, err := DeserializeTokenTimerResponseBody(harEntry)
	return resp, err
}

func (c *Client) DeleteTimers(reqCtx ApiRequestContext, ctxId string, tokId string) (*ApiResponse, error) {
	const semLogContext = "tpm-tokens-client::post-delete-timers"

	ep := c.tokenApiUrl(TokenTimersDelete, ctxId, tokId, nil)

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
