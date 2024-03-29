package tokensclient

import (
	"encoding/json"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/har"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-tokens-client/tokensclient/model/token"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
	"strings"
)

type TokenApiRequest struct {
	TokenId       string                 `yaml:"token-id,omitempty" mapstructure:"token-id,omitempty" json:"token-id,omitempty"`
	Typ           token.TokenType        `yaml:"type,omitempty" mapstructure:"type,omitempty" json:"type,omitempty"`
	CustomData    map[string]interface{} `yaml:"properties,omitempty" mapstructure:"properties,omitempty" json:"properties,omitempty"`
	CheckOnlyFLag bool                   `yaml:"check-only,omitempty" mapstructure:"check-only,omitempty" json:"check-only,omitempty"`
}

func (tok *TokenApiRequest) ToJSON() ([]byte, error) {
	return json.Marshal(tok)
}

func (c *Client) GetToken(reqCtx ApiRequestContext, ctxId string, tokId string) (*token.Token, error) {
	const semLogContext = "tpm-tokens-client::get-token"

	ep := c.tokenApiUrl(GetToken, ctxId, tokId, "", nil)

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
	// c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeTokenResponseBody(harEntry)
	return resp, err
}

func (c *Client) NewToken(reqCtx ApiRequestContext, ctxId string, tokenRequest *TokenApiRequest, ct string) (*token.Token, error) {
	const semLogContext = "tpm-tokens-client::new-token"

	op := "use"
	if tokenRequest.CheckOnlyFLag {
		op = "check"
	}

	ep := c.tokenApiUrl(NewToken, ctxId, "", "", []har.NameValuePair{{Name: "op", Value: op}})

	if ct == "" {
		ct = ContentTypeApplicationJson
	}

	switch ct {
	case ContentTypeApplicationJson:
	default:
		log.Warn().Str("content-type", ct).Msg(semLogContext + " unsupported content-type...using json")
		ct = ContentTypeApplicationJson
	}

	tokenRequest.TokenId = token.WellFormTokenId(tokenRequest.TokenId)
	b, err := tokenRequest.ToJSON()
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
	// c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeTokenResponseBody(harEntry)
	return resp, err
}

func (c *Client) DeleteToken(reqCtx ApiRequestContext, ctxId string, tokId string) (*ApiResponse, error) {
	const semLogContext = "tpm-tokens-client::delete-token"

	ep := c.tokenApiUrl(DeleteToken, ctxId, tokId, "", nil)

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
	// c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeApiResponse(harEntry)
	return resp, err
}

func (c *Client) CommitToken(reqCtx ApiRequestContext, ctxId string, tokId string) (*token.Token, error) {
	const semLogContext = "tpm-tokens-client::commit-token"

	ep := c.tokenApiUrl(TokenCommit, ctxId, tokId, "", nil)

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
	// c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeTokenResponseBody(harEntry)
	return resp, err
}

func (c *Client) RollbackToken(reqCtx ApiRequestContext, ctxId string, tokId string) (*token.Token, error) {
	const semLogContext = "tpm-tokens-client::commit-token"

	ep := c.tokenApiUrl(TokenRollback, ctxId, tokId, "", nil)

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
	// c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeTokenResponseBody(harEntry)
	return resp, err
}

func (c *Client) TokenNext(reqCtx ApiRequestContext, ctxId string, tokId string, tokenRequest *TokenApiRequest, ct string) (*token.Token, error) {
	return c.tokenNext(reqCtx, ctxId, tokId, "", tokenRequest, ct)
}

func (c *Client) TakeTransition(reqCtx ApiRequestContext, ctxId string, tokId string, transitionName string, tokenRequest *TokenApiRequest, ct string) (*token.Token, error) {
	return c.tokenNext(reqCtx, ctxId, tokId, transitionName, tokenRequest, ct)
}

// tokenNext private method to handle both the actual next op and the takeTransition one.
func (c *Client) tokenNext(reqCtx ApiRequestContext, ctxId string, tokId string, transitionName string, tokenRequest *TokenApiRequest, ct string) (*token.Token, error) {
	const semLogContext = "tpm-tokens-client::token-next"

	var ep string
	if transitionName == "" {
		ep = c.tokenApiUrl(TokenNext, ctxId, tokId, "", nil)
	} else {
		ep = c.tokenApiUrl(TokenTakeTransition, ctxId, tokId, transitionName, nil)
	}

	if ct == "" {
		ct = ContentTypeApplicationJson
	}

	switch ct {
	case ContentTypeApplicationJson:
	default:
		log.Warn().Str("content-type", ct).Msg(semLogContext + " unsupported content-type...using json")
		ct = ContentTypeApplicationJson
	}

	tokenRequest.TokenId = token.WellFormTokenId(tokenRequest.TokenId)
	b, err := tokenRequest.ToJSON()
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
	// c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeTokenResponseBody(harEntry)
	return resp, err
}

func (c *Client) TokenCheck(reqCtx ApiRequestContext, ctxId string, tokId string, tokenRequest *TokenApiRequest, ct string) (*token.Token, error) {
	const semLogContext = "tpm-tokens-client::token-check"

	ep := c.tokenApiUrl(TokenCheck, ctxId, tokId, "", nil)

	if ct == "" {
		ct = ContentTypeApplicationJson
	}

	switch ct {
	case ContentTypeApplicationJson:
	default:
		log.Warn().Str("content-type", ct).Msg(semLogContext + " unsupported content-type...using json")
		ct = ContentTypeApplicationJson
	}

	tokenRequest.TokenId = token.WellFormTokenId(tokenRequest.TokenId)
	b, err := tokenRequest.ToJSON()
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
	// c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeTokenResponseBody(harEntry)
	return resp, err
}

func (c *Client) tokenApiUrl(apiPath string, ctxId string, tokenId string, transitionName string, qParams []har.NameValuePair) string {
	var sb = strings.Builder{}
	sb.WriteString(c.host.Scheme)
	sb.WriteString("://")
	sb.WriteString(c.host.HostName)
	sb.WriteString(":")
	sb.WriteString(fmt.Sprint(c.host.Port))

	apiPath = strings.Replace(apiPath, TokenContextIdPathPlaceHolder, token.WellFormTokenContextId(ctxId), 1)
	apiPath = strings.Replace(apiPath, TokenIdPathPlaceHolder, token.WellFormTokenId(tokenId), 1)
	apiPath = strings.Replace(apiPath, TransitionNamePathPlaceHolder, transitionName, 1)
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
