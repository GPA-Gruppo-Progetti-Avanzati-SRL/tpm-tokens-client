package tokensclient

import (
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/har"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-tokens-client/tokensclient/model/businessview"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-tokens-client/tokensclient/model/token"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) GetTokenView(reqCtx ApiRequestContext, tokId string) (*businessview.Token, error) {
	const semLogContext = "tpm-tokens-client::get-token-view"

	ep := c.tokenViewApiUrl(GetTokenView, tokId, nil)

	req, err := c.client.NewRequest(http.MethodGet, ep, nil, reqCtx.getHeaders(""), nil)
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	harEntry, err := c.client.Execute(req,
		restclient.ExecutionWithOpName("client-token-view-get"),
		restclient.ExecutionWithRequestId(reqCtx.RequestId),
		restclient.ExecutionWithLraId(reqCtx.LRAId),
		restclient.ExecutionWithSpan(reqCtx.Span),
		restclient.ExecutionWithHarSpan(reqCtx.HarSpan))
	// c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeTokenViewContentResponse(harEntry)
	return resp, err
}

func (c *Client) GetActorView(reqCtx ApiRequestContext, actorId string, fullView bool) (*businessview.Actor, error) {
	const semLogContext = "tpm-tokens-client::get-actor-view"

	var qp []har.NameValuePair
	if fullView {
		qp = append(qp, har.NameValuePair{Name: "fullview", Value: "Y"})
	}
	ep := c.actorViewApiUrl(GetActorView, actorId, qp)

	req, err := c.client.NewRequest(http.MethodGet, ep, nil, reqCtx.getHeaders(""), nil)
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	harEntry, err := c.client.Execute(req,
		restclient.ExecutionWithOpName("client-actor-view-get"),
		restclient.ExecutionWithRequestId(reqCtx.RequestId),
		restclient.ExecutionWithLraId(reqCtx.LRAId),
		restclient.ExecutionWithSpan(reqCtx.Span),
		restclient.ExecutionWithHarSpan(reqCtx.HarSpan))
	// c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeActorViewContentResponse(harEntry)
	return resp, err
}

func (c *Client) tokenViewApiUrl(apiPath string, tokenId string, qParams []har.NameValuePair) string {
	var sb = strings.Builder{}
	sb.WriteString(c.host.Scheme)
	sb.WriteString("://")
	sb.WriteString(c.host.HostName)
	sb.WriteString(":")
	sb.WriteString(fmt.Sprint(c.host.Port))

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

func (c *Client) actorViewApiUrl(apiPath string, actorId string, qParams []har.NameValuePair) string {
	var sb = strings.Builder{}
	sb.WriteString(c.host.Scheme)
	sb.WriteString("://")
	sb.WriteString(c.host.HostName)
	sb.WriteString(":")
	sb.WriteString(fmt.Sprint(c.host.Port))

	apiPath = strings.Replace(apiPath, ActorIdPathPlaceHolder, businessview.WellFormActorId(actorId), 1)
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
