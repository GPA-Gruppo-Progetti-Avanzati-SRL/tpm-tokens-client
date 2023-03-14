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

/*
 * NewId API
 */

const (
	BridgeNewIdEndpointId = "new-token"
)

type NewTokenIdRequest struct {
	TokenContextId string                 `mapstructure:"context-id,omitempty"  json:"context-id,omitempty" yaml:"context-id,omitempty"`
	Id             string                 `mapstructure:"id,omitempty"  json:"id,omitempty" yaml:"id,omitempty"`
	Unique         bool                   `mapstructure:"unique"  json:"unique" yaml:"unique"`
	CustomData     map[string]interface{} `mapstructure:"properties,omitempty"  json:"properties,omitempty" yaml:"properties,omitempty"`
}

type NewTokenIdResponse struct {
	Id           string `yaml:"id,omitempty" mapstructure:"id,omitempty" json:"id,omitempty"`
	CreationDate string `yaml:"creation-date,omitempty" mapstructure:"creation-date,omitempty" json:"creation-date,omitempty"`
}

func (c *Client) NewId(reqCtx ApiRequestContext, ctxId string, unique bool, act map[string]interface{}) (*NewTokenIdResponse, error) {
	const semLogContext = "bridge-client::new-id"

	urlPath := c.findEndpointPathById(BridgeNewIdEndpointId)
	if urlPath == "" {
		log.Error().Msg(semLogContext + " unresolved endpoint url path")
	}

	ep := c.BridgeApiNewTokenUrl(urlPath, ctxId, nil)
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
		restclient.ExecutionWithOpName("client-new-token-context"),
		restclient.ExecutionWithRequestId(reqCtx.RequestId),
		restclient.ExecutionWithSpan(reqCtx.Span),
		restclient.ExecutionWithHarSpan(reqCtx.HarSpan))
	// c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	resp, err := DeserializeNewTokenIdResponseBody(harEntry)
	return resp, err
}

func DeserializeNewTokenIdResponseBody(resp *har.Entry) (*NewTokenIdResponse, error) {

	const semLogContext = "bridge-client::new-id-deserialize-response"
	if resp == nil || resp.Response == nil || resp.Response.Content == nil || resp.Response.Content.Data == nil {
		err := errors.New("cannot deserialize null response")
		log.Error().Err(err).Msg(semLogContext)
		return nil, NewExecutableServerError(WithErrorMessage(err.Error()))
	}

	var resultObj *NewTokenIdResponse
	var err error
	switch resp.Response.Status {
	case http.StatusOK:
		resultObj = &NewTokenIdResponse{}
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

func (c *Client) BridgeApiNewTokenUrl(apiPath string, ctxId string, qParams []har.NameValuePair) string {
	var sb = strings.Builder{}
	sb.WriteString(c.host.Scheme)
	sb.WriteString("://")
	sb.WriteString(c.host.HostName)
	sb.WriteString(":")
	sb.WriteString(fmt.Sprint(c.host.Port))
	sb.WriteString(strings.Replace(apiPath, "{context-id}", ctxId, 1))

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
