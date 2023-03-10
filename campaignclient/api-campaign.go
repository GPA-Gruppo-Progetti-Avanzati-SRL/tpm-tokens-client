package campaignclient

import (
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/har"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) GetCampaignById(reqCtx ApiRequestContext, ctxId string) (*Campaign, error) {
	const semLogContext = "campaign-client::get-campaign"

	ep := c.campaignApiUrl(CampaignGet, ctxId, nil)

	req, err := c.client.NewRequest(http.MethodGet, ep, nil, reqCtx.getHeaders(""), nil)
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	harEntry, err := c.client.Execute(req,
		restclient.ExecutionWithOpName(semLogContext),
		restclient.ExecutionWithRequestId(reqCtx.RequestId),
		restclient.ExecutionWithLraId(reqCtx.LRAId),
		restclient.ExecutionWithSpan(reqCtx.Span),
		restclient.ExecutionWithHarSpan(reqCtx.HarSpan))
	// c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeCampaignContentResponse(harEntry)
	return resp, err
}

func (c *Client) NewCampaign(reqCtx ApiRequestContext, tokenCtx *Campaign, ct string) (*Campaign, error) {
	const semLogContext = "campaign-client::new"
	ep := c.campaignApiUrl(CampaignNew, "", nil)

	if ct == "" {
		ct = ContentTypeApplicationJson
	}

	switch ct {
	case ContentTypeApplicationJson:
	default:
		log.Warn().Str("content-type", ct).Msg(semLogContext + " unsupported content-type...using json")
		ct = ContentTypeApplicationJson
	}

	b, err := tokenCtx.ToJSON()
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	req, err := c.client.NewRequest(http.MethodPost, ep, b, reqCtx.getHeaders(ct), nil)
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	harEntry, err := c.client.Execute(req,
		restclient.ExecutionWithOpName(semLogContext),
		restclient.ExecutionWithRequestId(reqCtx.RequestId),
		restclient.ExecutionWithLraId(reqCtx.LRAId),
		restclient.ExecutionWithSpan(reqCtx.Span),
		restclient.ExecutionWithHarSpan(reqCtx.HarSpan))
	// c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeCampaignContentResponse(harEntry)
	return resp, err
}

func (c *Client) ReplaceCampaign(reqCtx ApiRequestContext, tokenCtx *Campaign, ct string) (*Campaign, error) {
	const semLogContext = "campaign-client::replace"
	ep := c.campaignApiUrl(CampaignPut, tokenCtx.Id, nil)

	if ct == "" {
		ct = ContentTypeApplicationJson
	}

	switch ct {
	case ContentTypeApplicationJson:
	default:
		log.Warn().Str("content-type", ct).Msg(semLogContext + " unsupported content-type...using json")
		ct = ContentTypeApplicationJson
	}

	b, err := tokenCtx.ToJSON()
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	req, err := c.client.NewRequest(http.MethodPut, ep, b, reqCtx.getHeaders(ct), nil)
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	harEntry, err := c.client.Execute(req,
		restclient.ExecutionWithOpName(semLogContext),
		restclient.ExecutionWithRequestId(reqCtx.RequestId),
		restclient.ExecutionWithLraId(reqCtx.LRAId),
		restclient.ExecutionWithSpan(reqCtx.Span),
		restclient.ExecutionWithHarSpan(reqCtx.HarSpan))
	// c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeCampaignContentResponse(harEntry)
	return resp, err
}

func (c *Client) DeleteCampaign(reqCtx ApiRequestContext, ctxId string) (bool, error) {
	const semLogContext = "campaign-client::delete"

	ep := c.campaignApiUrl(CampaignDelete, ctxId, nil)

	req, err := c.client.NewRequest(http.MethodDelete, ep, nil, reqCtx.getHeaders(""), nil)
	if err != nil {
		return false, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	harEntry, err := c.client.Execute(req,
		restclient.ExecutionWithOpName(semLogContext),
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

func (c *Client) campaignApiUrl(apiPath string, ctxId string, qParams []har.NameValuePair) string {
	var sb = strings.Builder{}
	sb.WriteString(c.host.Scheme)
	sb.WriteString("://")
	sb.WriteString(c.host.HostName)
	sb.WriteString(":")
	sb.WriteString(fmt.Sprint(c.host.Port))
	sb.WriteString(strings.Replace(apiPath, CampaignIdPathPlaceHolder, ctxId, 1))

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
