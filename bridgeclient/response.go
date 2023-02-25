package bridgeclient

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
	"time"
)

const (
	StatusOKDefaultMessage    = "success"
	ErrorDefaultMessage       = ""
	ServerErrorDefaultMessage = "server error"
	BadRequestDefaultMessage  = "bad request"
)

type Option func(executableError *ApiResponse)

type ApiResponse struct {
	StatusCode  int    `yaml:"-" mapstructure:"-" json:"-"`
	ErrCode     string `json:"error-code,omitempty" yaml:"error-code,omitempty" mapstructure:"error-code,omitempty"`
	Ambit       string `json:"ambit,omitempty" yaml:"ambit,omitempty" mapstructure:"ambit,omitempty"`
	Step        string `yaml:"step,omitempty" mapstructure:"step,omitempty" json:"step,omitempty"`
	Text        string `json:"text,omitempty" yaml:"text,omitempty" mapstructure:"text,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`
	Message     string `yaml:"message,omitempty" mapstructure:"message,omitempty" json:"message,omitempty"`
	Ts          string `yaml:"timestamp,omitempty" mapstructure:"timestamp,omitempty" json:"timestamp,omitempty"`
}

func (ae *ApiResponse) Error() string {
	var sv strings.Builder
	const sep = " - "
	if ae.StatusCode != 0 {
		sv.WriteString(fmt.Sprintf("status-code: %d"+sep, ae.StatusCode))
	}

	if ae.ErrCode != "" {
		sv.WriteString(fmt.Sprintf("error-code: %s"+sep, ae.ErrCode))
	}

	if ae.Ambit != "" {
		sv.WriteString(fmt.Sprintf("ambit: %s"+sep, ae.Ambit))
	}

	if ae.Step != "" {
		sv.WriteString(fmt.Sprintf("step: %s"+sep, ae.Step))
	}

	if ae.Text != "" {
		sv.WriteString(fmt.Sprintf("text: %s"+sep, ae.Text))
	}

	if ae.Description != "" {
		sv.WriteString(fmt.Sprintf("description: %s"+sep, ae.Description))
	}

	if ae.Message != "" {
		sv.WriteString(fmt.Sprintf("message: %s"+sep, ae.Message))
	}

	if ae.Ts != "" {
		sv.WriteString(fmt.Sprintf("timestamp: %s"+sep, ae.Ts))
	}

	return strings.TrimSuffix(sv.String(), sep)
}

func DeserApiResponseFromJson(b []byte) (ApiResponse, error) {
	a := ApiResponse{}
	err := json.Unmarshal(b, &a)
	return a, err
}

func WithErrorStatusCode(c int) Option {
	return func(e *ApiResponse) {
		e.StatusCode = c
	}
}

func WithCode(c string) Option {
	return func(e *ApiResponse) {
		e.ErrCode = c
	}
}

func WithErrorMessage(m string) Option {
	return func(e *ApiResponse) {
		e.Text = m
	}
}

func WithDescription(m string) Option {
	return func(e *ApiResponse) {
		e.Description = m
	}
}

func NewExecutableError(opts ...Option) *ApiResponse {
	err := &ApiResponse{StatusCode: 0, Text: ErrorDefaultMessage}
	for _, o := range opts {
		o(err)
	}
	return err
}

func NewExecutableServerError(opts ...Option) *ApiResponse {
	err := &ApiResponse{StatusCode: http.StatusInternalServerError, Text: ServerErrorDefaultMessage}
	for _, o := range opts {
		o(err)
	}
	return err
}

func NewSuccessResponse(opts ...Option) *ApiResponse {
	err := &ApiResponse{StatusCode: http.StatusOK, Text: StatusOKDefaultMessage}
	for _, o := range opts {
		o(err)
	}
	return err
}

func NewBadRequestError(opts ...Option) *ApiResponse {
	err := &ApiResponse{StatusCode: http.StatusBadRequest, Text: BadRequestDefaultMessage}
	for _, o := range opts {
		o(err)
	}
	return err
}

func (exe *ApiResponse) ToJSON() []byte {
	exe.Ts = time.Now().Format(time.RFC3339)
	var b []byte
	var err error

	b, err = json.Marshal(exe)
	if err != nil {
		log.Error().Err(err).Msg("error in marshalling api-error")
		return []byte(`{"msg": "error in marshalling api-error"}`)
	}

	return b
}

func ErrorCode(err error) string {
	if resp, ok := err.(*ApiResponse); ok {
		return resp.ErrCode
	}

	return ""
}
