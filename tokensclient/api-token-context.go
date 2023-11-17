package tokensclient

import (
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/har"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-tokens-client/tokensclient/model/token"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) GetTokenContextById(reqCtx ApiRequestContext, ctxId string) (*token.TokenContext, error) {
	const semLogContext = "tpm-tokens-client::get-token-context"

	ep := c.tokenContextApiUrl(TokenContextGet, ctxId, nil)

	req, err := c.client.NewRequest(http.MethodGet, ep, nil, reqCtx.getHeaders(""), nil)
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	harEntry, err := c.client.Execute(req,
		restclient.ExecutionWithOpName("client-get-token-context"),
		restclient.ExecutionWithRequestId(reqCtx.RequestId),
		restclient.ExecutionWithLraId(reqCtx.LRAId),
		restclient.ExecutionWithSpan(reqCtx.Span),
		restclient.ExecutionWithHarSpan(reqCtx.HarSpan))
	// c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeTokenContextContentResponse(harEntry)
	return resp, err
}

func (c *Client) NewTokenContext(reqCtx ApiRequestContext, tokenCtx *token.TokenContext, ct string) (*token.TokenContext, error) {
	const semLogContext = "tpm-tokens-client::new-token-context"

	ep := c.tokenContextApiUrl(TokenContextNew, "", nil)

	if ct == "" {
		ct = ContentTypeApplicationJson
	}

	switch ct {
	case ContentTypeApplicationJson:
	default:
		log.Warn().Str("content-type", ct).Msg(semLogContext + " unsupported content-type...using json")
		ct = ContentTypeApplicationJson
	}

	tokenCtx.Id = token.WellFormTokenContextId(tokenCtx.Id)
	b, err := tokenCtx.ToJSON()
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	req, err := c.client.NewRequest(http.MethodPost, ep, b, reqCtx.getHeaders(ct), nil)
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	harEntry, err := c.client.Execute(req,
		restclient.ExecutionWithOpName("client-new-token-context"),
		restclient.ExecutionWithRequestId(reqCtx.RequestId),
		restclient.ExecutionWithLraId(reqCtx.LRAId),
		restclient.ExecutionWithSpan(reqCtx.Span),
		restclient.ExecutionWithHarSpan(reqCtx.HarSpan))
	// c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeTokenContextContentResponse(harEntry)
	return resp, err
}

func (c *Client) ReplaceTokenContext(reqCtx ApiRequestContext, tokenCtx *token.TokenContext, ct string) (*token.TokenContext, error) {
	const semLogContext = "tpm-tokens-client::new-token-context"

	ep := c.tokenContextApiUrl(TokenContextPut, tokenCtx.Id, nil)

	if ct == "" {
		ct = ContentTypeApplicationJson
	}

	switch ct {
	case ContentTypeApplicationJson:
	default:
		log.Warn().Str("content-type", ct).Msg(semLogContext + " unsupported content-type...using json")
		ct = ContentTypeApplicationJson
	}

	tokenCtx.Id = token.WellFormTokenContextId(tokenCtx.Id)
	b, err := tokenCtx.ToJSON()
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	req, err := c.client.NewRequest(http.MethodPut, ep, b, reqCtx.getHeaders(ct), nil)
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	harEntry, err := c.client.Execute(req,
		restclient.ExecutionWithOpName("client-replace-token-context"),
		restclient.ExecutionWithRequestId(reqCtx.RequestId),
		restclient.ExecutionWithLraId(reqCtx.LRAId),
		restclient.ExecutionWithSpan(reqCtx.Span),
		restclient.ExecutionWithHarSpan(reqCtx.HarSpan))
	// c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeTokenContextContentResponse(harEntry)
	return resp, err
}

func (c *Client) DeleteTokenContext(reqCtx ApiRequestContext, ctxId string) (bool, error) {
	const semLogContext = "tpm-tokens-client::delete-token-context"

	ep := c.tokenContextApiUrl(TokenContextDelete, ctxId, nil)

	req, err := c.client.NewRequest(http.MethodDelete, ep, nil, reqCtx.getHeaders(""), nil)
	if err != nil {
		return false, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	harEntry, err := c.client.Execute(req,
		restclient.ExecutionWithOpName("client-delete-token-context"),
		restclient.ExecutionWithRequestId(reqCtx.RequestId),
		restclient.ExecutionWithLraId(reqCtx.LRAId),
		restclient.ExecutionWithSpan(reqCtx.Span),
		restclient.ExecutionWithHarSpan(reqCtx.HarSpan))
	// c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return false, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeApiResponse(harEntry)
	if err != nil {
		return false, err
	}

	rc := false
	switch resp.StatusCode {
	case http.StatusOK:
		rc = true
	case http.StatusNotFound:
	default:
		// Return an error if not ok or not entity not found.
		err = resp
	}

	return rc, err
}

func (c *Client) tokenContextApiUrl(apiPath string, ctxId string, qParams []har.NameValuePair) string {
	var sb = strings.Builder{}
	sb.WriteString(c.host.Scheme)
	sb.WriteString("://")
	sb.WriteString(c.host.HostName)
	sb.WriteString(":")
	sb.WriteString(fmt.Sprint(c.host.Port))
	sb.WriteString(strings.Replace(apiPath, TokenContextIdPathPlaceHolder, token.WellFormTokenContextId(ctxId), 1))

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
