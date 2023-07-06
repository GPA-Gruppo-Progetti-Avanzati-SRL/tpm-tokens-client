package token

import (
	"fmt"
	"net/http"
)

const (
	TokenErrorSystem                         = "tok-sys-err"
	TokenErrorSystemConfiguration            = "tok-cfg-err"
	TokenErrorExpressionEvaluation           = "tok-expr-err"
	TokenErrorPropertiesValidationEvaluation = "tok-properties-valid-err"
	TokenErrorContextDefinition              = "tok-ctx-def-err"
	TokenErrorNotTransitionFound             = "tok-no-transition-err"
	TokenErrorNewTokenId                     = "tok-new-tok-id-err"
	TokenErrorTransactionInvalidState        = "tok-tx-inv-state-err"
	TokenErrorInvalidState                   = "tok-inv-state-err"
	TokenExpiredError                        = "tok-expired-err"
	TokenFinalStateAlreadyReachedError       = "tok-final-state-err"
	TokenAlreadyExists                       = "tok-already-exists"
	TokenDupRequestError                     = "tok-err-dup-request"
	TokenNotFoundError                       = "tok-not-found-err"
	TokenErrorGenericText                    = "generic error"
	TokenContextNotActiveError               = "tok-ctx-not-active"
	TokenContextNotFoundError                = "tok-ctx-not-found"
	TokenContextAlreadyExists                = "tok-ctx-already-exists"
)

type TokErrorInfo struct {
	Code       string
	Text       string
	StatusCode int
}

var TokErrorTextMapping = map[string]TokErrorInfo{
	TokenErrorSystem:                         {StatusCode: http.StatusBadRequest, Code: TokenErrorSystem, Text: "general error"},
	TokenErrorSystemConfiguration:            {StatusCode: http.StatusBadRequest, Code: TokenErrorSystemConfiguration, Text: "system error"},
	TokenErrorExpressionEvaluation:           {StatusCode: http.StatusInternalServerError, Code: TokenErrorExpressionEvaluation, Text: "expression evaluation error"},
	TokenErrorPropertiesValidationEvaluation: {StatusCode: http.StatusPreconditionFailed, Code: TokenErrorPropertiesValidationEvaluation, Text: "input params validation"},
	TokenErrorContextDefinition:              {StatusCode: http.StatusInternalServerError, Code: TokenErrorContextDefinition, Text: "context definition error"},
	TokenErrorNotTransitionFound:             {StatusCode: http.StatusPreconditionFailed, Code: TokenErrorNotTransitionFound, Text: "transition not found"},
	TokenErrorNewTokenId:                     {StatusCode: http.StatusInternalServerError, Code: TokenErrorNewTokenId, Text: "new token id error"},
	TokenErrorTransactionInvalidState:        {StatusCode: http.StatusConflict, Code: TokenErrorTransactionInvalidState, Text: "invalid transactional state"},
	TokenErrorInvalidState:                   {StatusCode: http.StatusConflict, Code: TokenErrorInvalidState, Text: "invalid state"},
	TokenFinalStateAlreadyReachedError:       {StatusCode: http.StatusPreconditionFailed, Code: TokenFinalStateAlreadyReachedError, Text: "token final state already reached"},
	TokenDupRequestError:                     {StatusCode: http.StatusConflict, Code: TokenDupRequestError, Text: "the request has already been processed"},
	TokenContextNotFoundError:                {StatusCode: http.StatusBadRequest, Code: TokenContextNotFoundError, Text: "token context not found"},
	TokenContextAlreadyExists:                {StatusCode: http.StatusBadRequest, Code: TokenContextAlreadyExists, Text: "token context already exists"},
	TokenContextNotActiveError:               {StatusCode: http.StatusBadRequest, Code: TokenContextNotActiveError, Text: "token context not active"},
	TokenExpiredError:                        {StatusCode: http.StatusConflict, Code: TokenExpiredError, Text: "token expired"},
	TokenNotFoundError:                       {StatusCode: http.StatusNotFound, Code: TokenNotFoundError, Text: "token not found"},
}

type TokError struct {
	Code        string `yaml:"code,omitempty" mapstructure:"code,omitempty" json:"code,omitempty"`
	Text        string `yaml:"text,omitempty" mapstructure:"text,omitempty" json:"text,omitempty"`
	Description string `yaml:"description,omitempty" mapstructure:"description,omitempty" json:"description,omitempty"`
}

func (te *TokError) Error() string {
	return fmt.Sprintf("%s - %s", te.Code, te.Text)
}

func NewTokError(c string, d string) error {
	t := TokError{Code: c, Description: d, Text: MapErrorCode2TokErrorInfo(c).Text}
	return &t
}

func MapErrorCode2TokErrorInfo(c string) TokErrorInfo {

	ti := TokErrorInfo{Code: c, Text: TokenErrorGenericText, StatusCode: http.StatusBadRequest}
	if ti, ok := TokErrorTextMapping[c]; ok {
		return ti
	}

	return ti
}

/*
func MapErrorCode2Text(errCode string) string {
	t := TokenErrorGenericText
	if text, ok := TokErrorTextMapping[errCode]; ok {
		t = text.Text
	}

	return t
}
*/
