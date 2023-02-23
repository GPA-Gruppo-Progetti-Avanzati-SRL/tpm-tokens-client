package tokensclient

import "fmt"

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

var TokenErrorTextMapping = map[string]string{
	TokenErrorSystem:                   "general error",
	TokenErrorSystemConfiguration:      "system error",
	TokenErrorExpressionEvaluation:     "expression evaluation error",
	TokenErrorContextDefinition:        "context definition error",
	TokenErrorNotTransitionFound:       "transition not found",
	TokenErrorNewTokenId:               "new token id error",
	TokenErrorTransactionInvalidState:  "invalid transactional state",
	TokenErrorInvalidState:             "invalid state",
	TokenFinalStateAlreadyReachedError: "token final state already reached",
	TokenDupRequestError:               "the request has already been processed",
	TokenContextNotFoundError:          "token context not found",
	TokenContextAlreadyExists:          "token context already exists",
	TokenContextNotActiveError:         "token context not active",
	TokenExpiredError:                  "token expired",
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
	t := TokError{Code: c, Description: d, Text: MapErrorCode2Text(c)}
	return &t
}

func MapErrorCode2Text(errCode string) string {
	t := TokenErrorGenericText
	if text, ok := TokenErrorTextMapping[errCode]; ok {
		t = text
	}

	return t
}
