package tokensclient

import (
	"fmt"
	"net/http"
)

const (
	BerSystemError        = "ber-sys-err"
	BerNotFoundError      = "ber-not-found"
	BerAlreadyExistsError = "ber-already-exists"
	BerErrorGenericText   = "generic error"
)

type BerErrorInfo struct {
	Code       string
	Text       string
	StatusCode int
}

var BerErrorTextMapping = map[string]BerErrorInfo{
	BerSystemError:        {StatusCode: http.StatusBadRequest, Code: BerSystemError, Text: "general error"},
	BerNotFoundError:      {StatusCode: http.StatusBadRequest, Code: BerNotFoundError, Text: "bearer not found in context"},
	BerAlreadyExistsError: {StatusCode: http.StatusBadRequest, Code: BerAlreadyExistsError, Text: "bearer already enlisted in context"},
}

type BerError struct {
	Code        string `yaml:"code,omitempty" mapstructure:"code,omitempty" json:"code,omitempty"`
	Text        string `yaml:"text,omitempty" mapstructure:"text,omitempty" json:"text,omitempty"`
	Description string `yaml:"description,omitempty" mapstructure:"description,omitempty" json:"description,omitempty"`
}

func (te *BerError) Error() string {
	return fmt.Sprintf("%s - %s", te.Code, te.Text)
}

func NewBerError(c string, d string) error {
	t := BerError{Code: c, Description: d, Text: MapErrorCode2BerErrorInfo(c).Text}
	return &t
}

func MapErrorCode2BerErrorInfo(c string) BerErrorInfo {

	ti := BerErrorInfo{Code: c, Text: BerErrorGenericText, StatusCode: http.StatusBadRequest}
	if ti, ok := BerErrorTextMapping[c]; ok {
		return ti
	}

	return ti
}

/*
func MapErrorCode2Text(errCode string) string {
	t := BerErrorGenericText
	if text, ok := BerErrorTextMapping[errCode]; ok {
		t = text.Text
	}

	return t
}
*/
