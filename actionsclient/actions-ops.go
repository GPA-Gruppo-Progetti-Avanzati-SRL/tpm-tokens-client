package actionsclient

import (
	"encoding/json"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/expression"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/har"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func (lks *LinkedService) NewClientById(id string, expressionCtx *expression.Context, opts ...restclient.Option) (*Client, error) {

	var err error
	var actCfg Config
	for _, a := range lks.cfg {
		if a.Id == id {
			actCfg = a
		}
	}

	if actCfg.Id != id {
		return nil, fmt.Errorf("could not resolve action %s", id)
	}

	client, err := lks.NewClient(actCfg, expressionCtx, opts...)
	return client, err
}

func (lks *LinkedService) NewClient(cfg Config, expressionCtx *expression.Context, opts ...restclient.Option) (*Client, error) {
	const semLogContext = semLogContextBase + "::new"

	resolvedCfg := cfg
	if expressionCtx != nil {
		v, err := expressionCtx.EvalOne(cfg.Path)
		if err != nil {
			return nil, err
		}
		resolvedCfg.Path = fmt.Sprint(v)

		resolvedCfg.Headers = nil
		for i := range cfg.Headers {
			v, err := expressionCtx.EvalOne(cfg.Headers[i].Value)
			if err != nil {
				return nil, err
			}
			resolvedCfg.Headers = append(resolvedCfg.Headers, restclient.Header{Name: cfg.Headers[i].Name, Value: fmt.Sprint(v)})
		}
	}

	client := restclient.NewClient(&resolvedCfg.Config, opts...)

	h := cfg.Host.FixValues()
	log.Trace().Str("scheme", h.Scheme).Int("port", h.Port).Str("host-name", h.HostName).Msg(semLogContext)
	return &Client{client: client, host: h, method: resolvedCfg.Method, path: resolvedCfg.Path, useResponse: resolvedCfg.Type == ActionTypeEnrich}, nil
}

func (c *Client) ExecuteAction(actionId string, actionBody map[string]interface{}) (map[string]interface{}, error) {

	const semLogContext = semLogContextBase + "::execute-action"

	b, err := json.Marshal(actionBody)
	if err != nil {
		return nil, err
	}

	req, err := c.client.NewRequest(c.method, c.Url(nil), b, nil, nil)
	if err != nil {
		return nil, err
	}

	harEntry, err := c.client.Execute(req, restclient.ExecutionWithOpName("actions-client"), restclient.ExecutionWithRequestId("auto-req-id"))
	c.harEntries = append(c.harEntries, harEntry)
	if err != nil {
		return nil, &ActionResponse{
			StatusCode: harEntry.Response.Status,
			Message:    err.Error(),
			Ts:         time.Now().Format(time.RFC3339Nano),
		}
	}

	log.Info().Str("action-id", actionId).Int("status-code", harEntry.Response.Status).Msg(semLogContext)

	if c.useResponse {
		var m map[string]interface{}
		m, err = handleEnrichingResponse(harEntry)
		if err == nil {
			for n, v := range m {
				actionBody[n] = v
			}
			return actionBody, nil
		}
	} else {
		err = handleBooleanResponse(harEntry)
	}

	return nil, err
}

func handleEnrichingResponse(harEntry *har.Entry) (map[string]interface{}, error) {

	const semLogContext = semLogContextBase + "::handle-enrich-response"
	var ar ActionResponse
	var err error

	sc := http.StatusInternalServerError
	if harEntry.Response != nil {
		sc = harEntry.Response.Status
	}

	if sc == 200 {
		if harEntry != nil && harEntry.Response != nil && harEntry.Response.Content != nil && harEntry.Response.Content.Data != nil {
			m := make(map[string]interface{})
			err = json.Unmarshal(harEntry.Response.Content.Data, &m)
			if err == nil {
				return m, nil
			} else {
				err = fmt.Errorf("error unmarshalling response of content-type: %s", harEntry.Response.Content.MimeType)
			}
		} else {
			err = fmt.Errorf("no response received from success action %d", sc)
		}
	} else {
		if harEntry != nil && harEntry.Response != nil && harEntry.Response.Content != nil && harEntry.Response.Content.Data != nil {
			err = json.Unmarshal(harEntry.Response.Content.Data, &ar)
			if err == nil {
				return nil, &ar
			}
		} else {
			err = fmt.Errorf("no response received from in error action %d", sc)
		}
	}

	ar = ActionResponse{
		StatusCode: http.StatusInternalServerError,
		Message:    err.Error(),
		Ts:         time.Now().Format(time.RFC3339Nano),
	}

	return nil, &ar
}

func handleBooleanResponse(harEntry *har.Entry) error {

	const semLogContext = semLogContextBase + "::handle-bool-response"
	var ar ActionResponse
	var err error

	sc := http.StatusInternalServerError
	if harEntry.Response != nil {
		sc = harEntry.Response.Status
	}

	if sc == 200 {
		return nil
	}

	if harEntry != nil && harEntry.Response != nil && harEntry.Response.Content != nil && harEntry.Response.Content.Data != nil {
		err = json.Unmarshal(harEntry.Response.Content.Data, &ar)
		if err == nil {
			return &ar
		} else {
			err = fmt.Errorf("error unmarshalling response of content-type: %s", harEntry.Response.Content.MimeType)
		}
	} else {
		err = fmt.Errorf("no response received from in error action %d", sc)
	}

	ar = ActionResponse{
		StatusCode: http.StatusInternalServerError,
		Message:    err.Error(),
		Ts:         time.Now().Format(time.RFC3339Nano),
	}

	return &ar
}

func (lks *LinkedService) CallActions(acts []string, expressionCtx *expression.Context, body map[string]interface{}, opts ...restclient.Option) (map[string]interface{}, error) {

	const semLogContext = semLogContextBase + "::call-actions"

	actionBody := make(map[string]interface{})
	for n, v := range body {
		actionBody[n] = v
	}

	for _, actId := range acts {

		c, ok := lks.FindConfigByActionId(actId)
		if !ok {
			log.Error().Str("action-id", actId).Msg(semLogContext + " action not found")
			return nil, &ActionResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    fmt.Sprintf("action %s not found", actId),
				Ts:         time.Now().Format(time.RFC3339Nano),
			}
		} else {
			log.Trace().Str("action-id", actId).Msg(semLogContext)
		}

		cli, err := lks.NewClient(c, expressionCtx, opts...)
		if err != nil {
			return nil, err
		}

		m, err := cli.ExecuteAction(actId, actionBody)
		cli.Close()
		if err != nil {
			return nil, err
		}

		/*
			m, err := lks.executeAction(c, expressionCtx, actionBody, opts...)
			if err != nil {
				return nil, err
			}
		*/

		actionBody = m
	}

	return actionBody, nil
}

/*
func (lks *LinkedService) executeAction(c Config, expressionCtx *expression.Context, actionBody map[string]interface{}, opts ...restclient.Option) (map[string]interface{}, error) {

	const semLogContext = semLogContextBase + "::execute-action"
	cli, err := lks.NewClient(c, expressionCtx, opts...)
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	return cli.ExecuteAction(actionBody)
}
*/
