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

type TokenApiRequest struct {
	TokenId    string                 `yaml:"token-id,omitempty" mapstructure:"token-id,omitempty" json:"token-id,omitempty"`
	CustomData map[string]interface{} `yaml:"custom-data,omitempty" mapstructure:"custom-data,omitempty" json:"custom-data,omitempty"`
}

func (tok *TokenApiRequest) ToJSON() ([]byte, error) {
	return json.Marshal(tok)
}

func (c *Client) GetToken(reqCtx ApiRequestContext, ctxId string, tokId string) (*Token, error) {
	const semLogContext = "tpm-tokens::get-token"

	ep := c.tokenApiUrl(GetToken, ctxId, tokId, nil)

	req, err := c.client.NewRequest(http.MethodGet, ep, nil, reqCtx.getHeaders(""), nil)
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	harEntry, err := c.client.Execute(req,
		restclient.ExecutionWithOpName("client-token-get"),
		restclient.ExecutionWithRequestId(reqCtx.RequestId),
		restclient.ExecutionWithLraId(reqCtx.LRAId),
		restclient.ExecutionWithSpan(reqCtx.Span),
		restclient.ExecutionWithHarSpan(reqCtx.HarSpan))
	c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeTokenResponseBody(harEntry)
	return resp, err
}

func (c *Client) NewToken(reqCtx ApiRequestContext, ctxId string, token *TokenApiRequest, ct string) (*Token, error) {
	const semLogContext = "tpm-tokens::new-token"

	ep := c.tokenApiUrl(NewToken, ctxId, "", nil)

	if ct == "" {
		ct = ContentTypeApplicationJson
	}

	switch ct {
	case ContentTypeApplicationJson:
	default:
		log.Warn().Str("content-type", ct).Msg(semLogContext + " unsupported content-type...using json")
		ct = ContentTypeApplicationJson
	}

	b, err := token.ToJSON()
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	req, err := c.client.NewRequest(http.MethodPost, ep, b, reqCtx.getHeaders(ct), nil)
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	harEntry, err := c.client.Execute(req,
		restclient.ExecutionWithOpName("client-new-token"),
		restclient.ExecutionWithRequestId(reqCtx.RequestId),
		restclient.ExecutionWithLraId(reqCtx.LRAId),
		restclient.ExecutionWithSpan(reqCtx.Span),
		restclient.ExecutionWithHarSpan(reqCtx.HarSpan))
	c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeTokenResponseBody(harEntry)
	return resp, err
}

func (c *Client) DeleteToken(reqCtx ApiRequestContext, ctxId string, tokId string) (*ApiResponse, error) {
	const semLogContext = "tpm-tokens::delete-token"

	ep := c.tokenApiUrl(DeleteToken, ctxId, tokId, nil)

	req, err := c.client.NewRequest(http.MethodDelete, ep, nil, reqCtx.getHeaders(""), nil)
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	harEntry, err := c.client.Execute(req,
		restclient.ExecutionWithOpName("client-delete-token"),
		restclient.ExecutionWithRequestId(reqCtx.RequestId),
		restclient.ExecutionWithLraId(reqCtx.LRAId),
		restclient.ExecutionWithSpan(reqCtx.Span),
		restclient.ExecutionWithHarSpan(reqCtx.HarSpan))
	c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeApiResponse(harEntry)
	return resp, err
}

func (c *Client) CommitToken(reqCtx ApiRequestContext, ctxId string, tokId string) (*Token, error) {
	const semLogContext = "tpm-tokens::commit-token"

	ep := c.tokenApiUrl(TokenCommit, ctxId, tokId, nil)

	req, err := c.client.NewRequest(http.MethodPut, ep, nil, reqCtx.getHeaders(""), nil)
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	harEntry, err := c.client.Execute(req,
		restclient.ExecutionWithOpName("client-token-commit"),
		restclient.ExecutionWithRequestId(reqCtx.RequestId),
		restclient.ExecutionWithLraId(reqCtx.LRAId),
		restclient.ExecutionWithSpan(reqCtx.Span),
		restclient.ExecutionWithHarSpan(reqCtx.HarSpan))
	c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeTokenResponseBody(harEntry)
	return resp, err
}

func (c *Client) RollbackToken(reqCtx ApiRequestContext, ctxId string, tokId string) (*Token, error) {
	const semLogContext = "tpm-tokens::commit-token"

	ep := c.tokenApiUrl(TokenRollback, ctxId, tokId, nil)

	req, err := c.client.NewRequest(http.MethodPut, ep, nil, reqCtx.getHeaders(""), nil)
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	harEntry, err := c.client.Execute(req,
		restclient.ExecutionWithOpName("client-token-rollback"),
		restclient.ExecutionWithRequestId(reqCtx.RequestId),
		restclient.ExecutionWithLraId(reqCtx.LRAId),
		restclient.ExecutionWithSpan(reqCtx.Span),
		restclient.ExecutionWithHarSpan(reqCtx.HarSpan))
	c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeTokenResponseBody(harEntry)
	return resp, err
}

func (c *Client) TokenNext(reqCtx ApiRequestContext, ctxId string, tokId string, token *TokenApiRequest, ct string) (*Token, error) {
	const semLogContext = "tpm-tokens::token-next"

	ep := c.tokenApiUrl(TokenNext, ctxId, tokId, nil)

	if ct == "" {
		ct = ContentTypeApplicationJson
	}

	switch ct {
	case ContentTypeApplicationJson:
	default:
		log.Warn().Str("content-type", ct).Msg(semLogContext + " unsupported content-type...using json")
		ct = ContentTypeApplicationJson
	}

	b, err := token.ToJSON()
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	req, err := c.client.NewRequest(http.MethodPut, ep, b, reqCtx.getHeaders(ct), nil)
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	harEntry, err := c.client.Execute(req,
		restclient.ExecutionWithOpName("client-token-next"),
		restclient.ExecutionWithRequestId(reqCtx.RequestId),
		restclient.ExecutionWithLraId(reqCtx.LRAId),
		restclient.ExecutionWithSpan(reqCtx.Span),
		restclient.ExecutionWithHarSpan(reqCtx.HarSpan))
	c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeTokenResponseBody(harEntry)
	return resp, err
}

func (c *Client) TokenCheck(reqCtx ApiRequestContext, ctxId string, tokId string, token *TokenApiRequest, ct string) (*Token, error) {
	const semLogContext = "tpm-tokens::token-check"

	ep := c.tokenApiUrl(TokenCheck, ctxId, tokId, nil)

	if ct == "" {
		ct = ContentTypeApplicationJson
	}

	switch ct {
	case ContentTypeApplicationJson:
	default:
		log.Warn().Str("content-type", ct).Msg(semLogContext + " unsupported content-type...using json")
		ct = ContentTypeApplicationJson
	}

	b, err := token.ToJSON()
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	req, err := c.client.NewRequest(http.MethodPut, ep, b, reqCtx.getHeaders(ct), nil)
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	harEntry, err := c.client.Execute(req,
		restclient.ExecutionWithOpName("client-token-check"),
		restclient.ExecutionWithRequestId(reqCtx.RequestId),
		restclient.ExecutionWithLraId(reqCtx.LRAId),
		restclient.ExecutionWithSpan(reqCtx.Span),
		restclient.ExecutionWithHarSpan(reqCtx.HarSpan))
	c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeTokenResponseBody(harEntry)
	return resp, err
}

func (c *Client) tokenApiUrl(apiPath string, ctxId string, tokenId string, qParams []har.NameValuePair) string {
	var sb = strings.Builder{}
	sb.WriteString(c.host.Scheme)
	sb.WriteString("://")
	sb.WriteString(c.host.HostName)
	sb.WriteString(":")
	sb.WriteString(fmt.Sprint(c.host.Port))

	apiPath = strings.Replace(apiPath, TokenContextIdPathPlaceHolder, ctxId, 1)
	apiPath = strings.Replace(apiPath, TokenIdPathPlaceHolder, tokenId, 1)
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
