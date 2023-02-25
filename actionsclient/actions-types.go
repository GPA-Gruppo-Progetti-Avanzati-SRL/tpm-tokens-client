package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/har"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
	"github.com/rs/zerolog/log"
	"net/url"
	"strings"
)

const semLogContextBase = "actions-client"

type LinkedService struct {
	cfg []Config
}

func NewInstanceWithConfig(cfg []Config) (*LinkedService, error) {
	lks := &LinkedService{cfg: cfg}
	return lks, nil
}

func (lks *LinkedService) FindConfigByActionId(actId string) (Config, bool) {
	for _, c := range lks.cfg {
		if c.Id == actId {
			return c, true
		}
	}

	return Config{}, false
}

type HostInfo struct {
	Scheme   string `mapstructure:"scheme,omitempty" json:"scheme,omitempty" yaml:"scheme,omitempty"`
	HostName string `mapstructure:"name,omitempty" json:"name,omitempty" yaml:"name,omitempty"`
	Port     int    `mapstructure:"port,omitempty" json:"port,omitempty" yaml:"port,omitempty"`
}

func (hi HostInfo) FixValues() HostInfo {

	h := HostInfo{
		Scheme:   hi.Scheme,
		Port:     hi.Port,
		HostName: hi.HostName,
	}

	if h.Scheme == "" {
		h.Scheme = "http"
	}

	if h.HostName == "" {
		h.HostName = "localhost"
	}

	if h.Port == 0 {
		switch h.Scheme {
		case "http":
			h.Port = 80
		case "https":
			h.Port = 443
		default:
			log.Error().Str("scheme", h.Scheme).Msg("host-info invalid scheme...reverting to http...")
		}
	}

	return h
}

type ActionType string

const (
	ActionTypeBool   = "bool"
	ActionTypeEnrich = "enrich"
)

type Config struct {
	restclient.Config `mapstructure:",squash"  yaml:",inline"`
	Id                string     `mapstructure:"id,omitempty" json:"id,omitempty" yaml:"id,omitempty"`
	Type              ActionType `mapstructure:"type,omitempty" json:"type,omitempty" yaml:"type,omitempty"`
	Host              HostInfo   `mapstructure:"host,omitempty" json:"host,omitempty" yaml:"host,omitempty"`
	Method            string     `mapstructure:"method,omitempty" json:"method,omitempty" yaml:"method,omitempty"`
	Path              string     `mapstructure:"path,omitempty" json:"path,omitempty" yaml:"path,omitempty"`
}

type Client struct {
	method      string
	path        string
	host        HostInfo
	client      *restclient.Client
	harEntries  []*har.Entry
	useResponse bool
}

func (c *Client) Close() {
	c.client.Close()
}

func (c *Client) Url(qParams []har.NameValuePair) string {
	var sb = strings.Builder{}
	sb.WriteString(c.host.Scheme)
	sb.WriteString("://")
	sb.WriteString(c.host.HostName)
	sb.WriteString(":")
	sb.WriteString(fmt.Sprint(c.host.Port))
	sb.WriteString(c.path)

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

type ActionResponse struct {
	StatusCode int    `yaml:"-" mapstructure:"-" json:"-"`
	ErrCode    string `json:"error-code,omitempty" yaml:"error-code,omitempty" mapstructure:"error-code,omitempty"`
	Text       string `json:"text,omitempty" yaml:"text,omitempty" mapstructure:"text,omitempty"`
	Message    string `yaml:"message,omitempty" mapstructure:"message,omitempty" json:"message,omitempty"`
	Ts         string `yaml:"timestamp,omitempty" mapstructure:"timestamp,omitempty" json:"timestamp,omitempty"`
}

func (ae *ActionResponse) Error() string {
	var sv strings.Builder
	const sep = " - "
	if ae.StatusCode != 0 {
		sv.WriteString(fmt.Sprintf("status-code: %d"+sep, ae.StatusCode))
	}

	if ae.ErrCode != "" {
		sv.WriteString(fmt.Sprintf("error-code: %s"+sep, ae.ErrCode))
	}

	if ae.Text != "" {
		sv.WriteString(fmt.Sprintf("text: %s"+sep, ae.Text))
	}

	if ae.Message != "" {
		sv.WriteString(fmt.Sprintf("message: %s"+sep, ae.Message))
	}

	if ae.Ts != "" {
		sv.WriteString(fmt.Sprintf("timestamp: %s"+sep, ae.Ts))
	}

	return strings.TrimSuffix(sv.String(), sep)
}

func DeserializeActionResponse(resp *har.Entry) (*ActionResponse, error) {
	const semLogContext = semLogContextBase + "::deserialize-action-response"
	if resp == nil || resp.Response == nil || resp.Response.Content == nil || resp.Response.Content.Data == nil {
		err := errors.New("cannot deserialize null response")
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	a := ActionResponse{}
	err := json.Unmarshal(resp.Response.Content.Data, &a)

	return &a, err
}
