package bridgeclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/har"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
	"strings"
)

const (
	RetrieveTokenEndpointId = "retrieve-token"
)

type RetrieveTokenRequest struct {
	TokenContextId string                 `mapstructure:"context-id,omitempty"  json:"context-id,omitempty" yaml:"context-id,omitempty"`
	TokenId        string                 `mapstructure:"token-id,omitempty"  json:"token-id,omitempty" yaml:"token-id,omitempty"`
	Unique         bool                   `mapstructure:"unique"  json:"unique" yaml:"unique"`
	Properties     map[string]interface{} `mapstructure:"properties,omitempty"  json:"properties,omitempty" yaml:"custom,omitempty"`
}

func (req *RetrieveTokenRequest) IsValid() bool {
	return true
}

type RetrieveTokenResponse struct {
	Id           string `yaml:"token-id,omitempty" mapstructure:"token-id,omitempty" json:"token-id,omitempty"`
	CreationDate string `yaml:"creation-date,omitempty" mapstructure:"creation-date,omitempty" json:"creation-date,omitempty"`
}

func (c *Client) RetrieveToken(reqCtx ApiRequestContext, ctxId string, tokenId string, unique bool, act map[string]interface{}) (*RetrieveTokenResponse, error) {
	const semLogContext = "bridge-client::retrieve-token"

	urlPath := c.findEndpointPathById(RetrieveTokenEndpointId)
	if urlPath == "" {
		log.Error().Msg(semLogContext + " unresolved endpoint url path")
	}

	ep := c.RetrieveTokenUrl(urlPath, ctxId, tokenId, nil)
	ct := ContentTypeApplicationJson

	b, err := json.Marshal(act)
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	req, err := c.client.NewRequest(http.MethodPost, ep, b, reqCtx.getHeaders(ct), nil)
	if err != nil {
		return nil, NewBadRequestError(WithErrorMessage(err.Error()))
	}

	harEntry, err := c.client.Execute(req,
		restclient.ExecutionWithOpName(RetrieveTokenEndpointId),
		restclient.ExecutionWithRequestId(reqCtx.RequestId),
		restclient.ExecutionWithSpan(reqCtx.Span),
		restclient.ExecutionWithHarSpan(reqCtx.HarSpan))
	// c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeRetrieveTokenResponseBody(harEntry)
	return resp, err
}

func DeserializeRetrieveTokenResponseBody(resp *har.Entry) (*RetrieveTokenResponse, error) {

	const semLogContext = "bridge-client::retrieve-token-deserialize-response"
	if resp == nil || resp.Response == nil || resp.Response.Content == nil || resp.Response.Content.Data == nil {
		err := errors.New("cannot deserialize null response")
		log.Error().Err(err).Msg(semLogContext)
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	var resultObj *RetrieveTokenResponse
	var err error
	switch resp.Response.Status {
	case http.StatusOK:
		resultObj = &RetrieveTokenResponse{}
		err = json.Unmarshal(resp.Response.Content.Data, resultObj)
		if err != nil {
			return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
		}

	default:
		var apiResponse ApiResponse
		apiResponse, err = DeserApiResponseFromJson(resp.Response.Content.Data)
		if err != nil {
			return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
		}
		apiResponse.StatusCode = resp.Response.Status
		err = &apiResponse
		return nil, err
	}

	return resultObj, nil
}

func (c *Client) RetrieveTokenUrl(apiPath string, ctxId string, tokenId string, qParams []har.NameValuePair) string {
	var sb = strings.Builder{}
	sb.WriteString(c.host.Scheme)
	sb.WriteString("://")
	sb.WriteString(c.host.HostName)
	sb.WriteString(":")
	sb.WriteString(fmt.Sprint(c.host.Port))
	apiPath = strings.Replace(apiPath, "{context-id}", ctxId, 1)
	sb.WriteString(strings.Replace(apiPath, "{token-id}", tokenId, 1))

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
