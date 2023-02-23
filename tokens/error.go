package tokens

import (
	"fmt"
	"net/http"
)

const (
	TokenErrorSystem                   = "tok-err-001"
	TokenErrorSystemConfiguration      = "tok-err-002"
	TokenErrorExpressionEvaluation     = "tok-err-003"
	TokenErrorContextDefinition        = "tok-err-004"
	TokenErrorNotTransitionFound       = "tok-err-005"
	TokenErrorNewTokenId               = "tok-err-006"
	TokenErrorTransactionInvalidState  = "tok-err-007"
	TokenErrorInvalidState             = "tok-err-008"
	TokenExpiredError                  = "tok-err-009"
	TokenFinalStateAlreadyReachedError = "tok-err-010"
	TokenDupRequestError               = "tok-err-dup-request"
	TokenErrorGenericText              = "generic error"
	TokenContextNotActiveError         = "tok-ctx-not-active"
	TokenContextNotFoundError          = "tok-ctx-not-found"
	TokenContextAlreadyExists          = "tok-ctx-already-exists"
)

type TokenErrorInfo struct {
	Code       string
	Text       string
	StatusCode int
}

var TokenErrorTextMapping = map[string]TokenErrorInfo{
	TokenErrorSystem:                   {StatusCode: http.StatusBadRequest, Code: TokenErrorSystem, Text: "general error"},
	TokenErrorSystemConfiguration:      {StatusCode: http.StatusBadRequest, Code: TokenErrorSystemConfiguration, Text: "system error"},
	TokenErrorExpressionEvaluation:     {StatusCode: http.StatusInternalServerError, Code: TokenErrorExpressionEvaluation, Text: "expression evaluation error"},
	TokenErrorContextDefinition:        {StatusCode: http.StatusInternalServerError, Code: TokenErrorContextDefinition, Text: "context definition error"},
	TokenErrorNotTransitionFound:       {StatusCode: http.StatusBadRequest, Code: TokenErrorNotTransitionFound, Text: "transition not found"},
	TokenErrorNewTokenId:               {StatusCode: http.StatusInternalServerError, Code: TokenErrorNewTokenId, Text: "new token id error"},
	TokenErrorTransactionInvalidState:  {StatusCode: http.StatusConflict, Code: TokenErrorTransactionInvalidState, Text: "invalid transactional state"},
	TokenErrorInvalidState:             {StatusCode: http.StatusConflict, Code: TokenErrorInvalidState, Text: "invalid state"},
	TokenFinalStateAlreadyReachedError: {StatusCode: http.StatusConflict, Code: TokenFinalStateAlreadyReachedError, Text: "token final state already reached"},
	TokenDupRequestError:               {StatusCode: http.StatusConflict, Code: TokenDupRequestError, Text: "the request has already been processed"},
	TokenContextNotFoundError:          {StatusCode: http.StatusBadRequest, Code: TokenContextNotFoundError, Text: "token context not found"},
	TokenContextAlreadyExists:          {StatusCode: http.StatusBadRequest, Code: TokenContextAlreadyExists, Text: "token context already exists"},
	TokenContextNotActiveError:         {StatusCode: http.StatusBadRequest, Code: TokenContextNotActiveError, Text: "token context not active"},
	TokenExpiredError:                  {StatusCode: http.StatusConflict, Code: TokenExpiredError, Text: "token expired"},
}

type TokError struct {
	Code        string `yaml:"code,omitempty" mapstructure:"code,omitempty" json:"code,omitempty"`
	Text        string `yaml:"text,omitempty" mapstructure:"text,omitempty" json:"text,omitempty"`
	Description string `yaml:"description,omitempty" mapstructure:"description,omitempty" json:"description,omitempty"`
}

func (te *TokError) Error() string {
	return fmt.Sprintf("%s - %s", te.Code, te.Text)
}

func NewError(c string, d string) error {
	t := TokError{Code: c, Description: d, Text: MapErrorCode(c).Text}
	return &t
}

func MapErrorCode(c string) TokenErrorInfo {

	ti := TokenErrorInfo{Code: c, Text: TokenErrorGenericText, StatusCode: http.StatusBadRequest}
	if ti, ok := TokenErrorTextMapping[c]; ok {
		return ti
	}

	return ti
}

/*
func MapErrorCode2Text(errCode string) string {
	t := TokenErrorGenericText
	if text, ok := TokenErrorTextMapping[errCode]; ok {
		t = text.Text
	}

	return t
}
*/
