package tokensclient

import "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"

const (
	BearerIdPathPlaceHolder       = "{bearer-id}"
	TokenContextIdPathPlaceHolder = "{context-id}"
	TokenIdPathPlaceHolder        = "{token-id}"

	TokenContextBasePath = "/api/v1/token-contexts"
	TokenContextQuery    = TokenContextBasePath
	TokenContextNew      = TokenContextBasePath
	TokenContextGet      = TokenContextBasePath + "/" + TokenContextIdPathPlaceHolder
	TokenContextPut      = TokenContextBasePath + "/" + TokenContextIdPathPlaceHolder
	TokenContextDelete   = TokenContextBasePath + "/" + TokenContextIdPathPlaceHolder

	TokenBasePath = TokenContextBasePath + "/" + TokenContextIdPathPlaceHolder + "/tokens"
	NewToken      = TokenBasePath
	GetToken      = TokenBasePath + "/" + TokenIdPathPlaceHolder
	DeleteToken   = TokenBasePath + "/" + TokenIdPathPlaceHolder
	TokenNext     = TokenBasePath + "/" + TokenIdPathPlaceHolder + "/next"
	TokenCheck    = TokenBasePath + "/" + TokenIdPathPlaceHolder + "/check"
	TokenCommit   = TokenBasePath + "/" + TokenIdPathPlaceHolder + "/commit"
	TokenRollback = TokenBasePath + "/" + TokenIdPathPlaceHolder + "/rollback"

	BearerBasePath                       = "/api/v1/bearers"
	BearerContextGet                     = BearerBasePath + "/" + BearerIdPathPlaceHolder + "/" + TokenContextIdPathPlaceHolder
	BearerContextPost                    = BearerBasePath + "/" + BearerIdPathPlaceHolder + "/" + TokenContextIdPathPlaceHolder
	BearerContextPut                     = BearerBasePath + "/" + BearerIdPathPlaceHolder + "/" + TokenContextIdPathPlaceHolder
	BearerContextDelete                  = BearerBasePath + "/" + BearerIdPathPlaceHolder + "/" + TokenContextIdPathPlaceHolder
	AddToken2BearerInContextPost         = BearerBasePath + "/" + BearerIdPathPlaceHolder + "/" + TokenContextIdPathPlaceHolder + "/" + TokenIdPathPlaceHolder
	RemoveTokenFromBearerInContextDelete = BearerBasePath + "/" + BearerIdPathPlaceHolder + "/" + TokenContextIdPathPlaceHolder + "/" + TokenIdPathPlaceHolder
)

type HostInfo struct {
	Scheme   string `mapstructure:"scheme,omitempty" json:"scheme,omitempty" yaml:"scheme,omitempty"`
	HostName string `mapstructure:"name,omitempty" json:"name,omitempty" yaml:"name,omitempty"`
	Port     int    `mapstructure:"port,omitempty" json:"port,omitempty" yaml:"port,omitempty"`
}

// Config Note: the json serialization seems not need any inline, squash of sorts...
type Config struct {
	restclient.Config `mapstructure:",squash"  yaml:",inline"`
	Host              HostInfo `mapstructure:"host,omitempty" json:"host,omitempty" yaml:"host,omitempty"`
}

func (c *Config) PostProcess() error {
	return nil
}
