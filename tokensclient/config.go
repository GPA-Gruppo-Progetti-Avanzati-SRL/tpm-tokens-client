package tokensclient

import "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"

const (
	ActorIdPathPlaceHolder        = "{actor-id}"
	TokenContextIdPathPlaceHolder = "{context-id}"
	TokenIdPathPlaceHolder        = "{token-id}"
	TransitionNamePathPlaceHolder = "{transition-name}"
	FactClassPathPlaceHolder      = "{fact-class}"
	FactGroupPathPlaceHolder      = "{fact-group}"
	FactIdPathPlaceHolder         = "{fact-id}"

	TokenContextBasePath = "/api/v1/token-contexts"
	TokenContextQuery    = TokenContextBasePath
	TokenContextNew      = TokenContextBasePath
	TokenContextGet      = TokenContextBasePath + "/" + TokenContextIdPathPlaceHolder
	TokenContextPut      = TokenContextBasePath + "/" + TokenContextIdPathPlaceHolder
	TokenContextDelete   = TokenContextBasePath + "/" + TokenContextIdPathPlaceHolder

	TokenBasePath       = TokenContextBasePath + "/" + TokenContextIdPathPlaceHolder + "/tokens"
	NewToken            = TokenBasePath
	GetToken            = TokenBasePath + "/" + TokenIdPathPlaceHolder
	DeleteToken         = TokenBasePath + "/" + TokenIdPathPlaceHolder
	TokenNext           = TokenBasePath + "/" + TokenIdPathPlaceHolder + "/next"
	TokenCheck          = TokenBasePath + "/" + TokenIdPathPlaceHolder + "/check"
	TokenCommit         = TokenBasePath + "/" + TokenIdPathPlaceHolder + "/commit"
	TokenRollback       = TokenBasePath + "/" + TokenIdPathPlaceHolder + "/rollback"
	TokenTimerCreate    = TokenBasePath + "/" + TokenIdPathPlaceHolder + "/timers"
	TokenTimersDelete   = TokenBasePath + "/" + TokenIdPathPlaceHolder + "/timers"
	TokenTakeTransition = TokenBasePath + "/" + TokenIdPathPlaceHolder + "/take/" + TransitionNamePathPlaceHolder

	BearerBasePath                       = "/api/v1/bearers"
	BearerContextGet                     = BearerBasePath + "/" + ActorIdPathPlaceHolder + "/" + TokenContextIdPathPlaceHolder
	BearerContextPost                    = BearerBasePath + "/" + ActorIdPathPlaceHolder + "/" + TokenContextIdPathPlaceHolder
	BearerContextPut                     = BearerBasePath + "/" + ActorIdPathPlaceHolder + "/" + TokenContextIdPathPlaceHolder
	BearerContextDelete                  = BearerBasePath + "/" + ActorIdPathPlaceHolder + "/" + TokenContextIdPathPlaceHolder
	AddToken2BearerInContextPost         = BearerBasePath + "/" + ActorIdPathPlaceHolder + "/" + TokenContextIdPathPlaceHolder + "/" + TokenIdPathPlaceHolder
	RemoveTokenFromBearerInContextDelete = BearerBasePath + "/" + ActorIdPathPlaceHolder + "/" + TokenContextIdPathPlaceHolder + "/" + TokenIdPathPlaceHolder

	ApiViewBasePath = "/api/v1/views"
	GetTokenView    = ApiViewBasePath + "/tokens/" + TokenIdPathPlaceHolder
	GetActorView    = ApiViewBasePath + "/actors/" + ActorIdPathPlaceHolder

	ApiFactsBasePath = "/api/v1/facts/"
	FactsQueryGroup  = ApiFactsBasePath + FactClassPathPlaceHolder + "/" + FactGroupPathPlaceHolder
	FactGet          = ApiFactsBasePath + FactClassPathPlaceHolder + "/" + FactGroupPathPlaceHolder + "/" + FactIdPathPlaceHolder
	FactAdd2Group    = ApiFactsBasePath + FactClassPathPlaceHolder + "/" + FactGroupPathPlaceHolder
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
