package campaignclient

import "fmt"

const (
	MalformedTokenError = "campaign-token-malformed"
	GenericErrorText    = "generic error"
	NotActiveError      = "campaign-not-active"
	NotFoundError       = "campaign-not-found"
	AlreadyExists       = "campaign-already-exists"
)

var ErrorTextMapping = map[string]string{
	NotFoundError:  "campaign not found",
	AlreadyExists:  "campaign already exists",
	NotActiveError: "campaign not active",
}

type Error struct {
	Code        string `yaml:"code,omitempty" mapstructure:"code,omitempty" json:"code,omitempty"`
	Text        string `yaml:"text,omitempty" mapstructure:"text,omitempty" json:"text,omitempty"`
	Description string `yaml:"description,omitempty" mapstructure:"description,omitempty" json:"description,omitempty"`
}

func (te *Error) Error() string {
	return fmt.Sprintf("%s - %s", te.Code, te.Text)
}

func NewError(c string, d string) error {
	t := Error{Code: c, Description: d, Text: MapErrorCode2Text(c)}
	return &t
}

func MapErrorCode2Text(errCode string) string {
	t := GenericErrorText
	if text, ok := ErrorTextMapping[errCode]; ok {
		t = text
	}

	return t
}
