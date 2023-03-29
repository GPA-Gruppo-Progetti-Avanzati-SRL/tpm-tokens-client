package tokensclient

import (
	"encoding/json"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/expression"
	"github.com/rs/zerolog/log"
	"time"
)

const (
	ContextPartitionKey = "token-context"

	TokenIdProviderTypeExternal    = "external"
	TokenIdProviderTypeDefault     = "default"
	TokenIdProviderTypeCosSequence = "cos-seq"

	ExpirationModeTimestamp = "timestamp"
	ExpirationModeDate      = "date"

	TokenContextBaseVersion = "v1.0.0"
)

type Timeline struct {
	StartDate      string `yaml:"start-date,omitempty" mapstructure:"start-date,omitempty" json:"start-date,omitempty"`
	EndDate        string `yaml:"end-date,omitempty" mapstructure:"end-date,omitempty" json:"end-date,omitempty"`
	ExpirationMode string `yaml:"expiration-mode,omitempty" mapstructure:"expiration-mode,omitempty" json:"expiration-mode,omitempty"`
}

func (tl *Timeline) IsOver() bool {
	today := time.Now().Format("20060102")
	if today > tl.EndDate {
		return true
	}

	return false
}

func (tl *Timeline) IsNotStartedYet() bool {
	today := time.Now().Format("20060102")
	if today < tl.StartDate {
		return true
	}

	return false
}

func (tl *Timeline) IsInRange() bool {
	today := time.Now().Format("20060102")
	if today >= tl.StartDate && today <= tl.EndDate {
		return true
	}

	return false
}

func (tl *Timeline) NumberOfDays() (int, error) {
	end, err := time.Parse("20060102", tl.EndDate)
	if err != nil {
		return -1, err
	}

	start, err := time.Parse("20060102", tl.StartDate)
	if err != nil {
		return -1, err
	}

	numDays := end.Sub(start).Hours() / 24.0
	return int(numDays), nil
}

func (tl *Timeline) Valid() bool {

	const semLogContext = "context timeline validation"

	if tl.StartDate == "" {
		tl.StartDate = time.Now().Format("20060102")
	} else {
		_, err := time.Parse("20060102", tl.StartDate)
		if err != nil {
			log.Error().Err(err).Str("start-date", tl.StartDate).Str("end-date", tl.EndDate).Msg(semLogContext + " invalid start-date format")
			return false
		}
	}

	if tl.EndDate == "" {
		tl.EndDate = time.Now().Format("20060102")
	} else {
		_, err := time.Parse("20060102", tl.EndDate)
		if err != nil {
			log.Error().Err(err).Str("start-date", tl.StartDate).Str("end-date", tl.EndDate).Msg(semLogContext + " invalid end-date format")
			return false
		}
	}

	if tl.StartDate > tl.EndDate {
		log.Error().Str("start-date", tl.StartDate).Str("end-date", tl.EndDate).Msg(semLogContext + " invalid range interval")
		return false
	}

	return true
}

type TokenIdProviderType struct {
	ProviderType string `yaml:"type,omitempty" mapstructure:"type,omitempty" json:"type,omitempty"`
	Unique       bool   `yaml:"unique,omitempty" mapstructure:"unique,omitempty" json:"unique,omitempty"`
	Format       string `yaml:"format,omitempty" mapstructure:"format,omitempty" json:"format,omitempty"`
}

type TokenContext struct {
	Id                  string               `yaml:"id,omitempty" mapstructure:"id,omitempty" json:"id,omitempty"`
	Pkey                string               `yaml:"pkey,omitempty" mapstructure:"pkey,omitempty" json:"pkey,omitempty"`
	Platform            string               `yaml:"platform,omitempty" mapstructure:"platform,omitempty" json:"platform,omitempty"`
	Version             string               `yaml:"version,omitempty" mapstructure:"version,omitempty" json:"version,omitempty"`
	Suspended           bool                 `yaml:"suspended,omitempty" mapstructure:"suspended,omitempty" json:"suspended,omitempty"`
	Timeline            Timeline             `yaml:"timeline,omitempty" mapstructure:"timeline,omitempty" json:"timeline,omitempty"`
	StateMachine        StateMachine         `yaml:"state-machine,omitempty" mapstructure:"state-machine,omitempty" json:"state-machine,omitempty"`
	TokenIdProviderType *TokenIdProviderType `yaml:"token-id-provider-type,omitempty" mapstructure:"token-id-provider-type,omitempty" json:"token-id-provider-type,omitempty"`
	TTL                 int                  `yaml:"ttl,omitempty" mapstructure:"ttl,omitempty" json:"ttl,omitempty"`
}

func (ctx *TokenContext) ToJSON() ([]byte, error) {
	return json.Marshal(ctx)
}

func (ctx *TokenContext) MustToJSON() []byte {
	b, err := json.Marshal(ctx)
	if err != nil {
		panic(err)
	}

	return b
}

func DeserializeContext(b []byte) (*TokenContext, error) {
	ctx := TokenContext{}
	err := json.Unmarshal(b, &ctx)
	if err != nil {
		return nil, err
	}

	ctx.Pkey = ContextPartitionKey
	return &ctx, nil
}

func (ctx *TokenContext) PostProcess() {

	// Initialize transition descriptions if they are empty.
	for i := range ctx.StateMachine.States {
		for j := range ctx.StateMachine.States[i].OutTransitions {
			if ctx.StateMachine.States[i].OutTransitions[j].Description == "" {
				ctx.StateMachine.States[i].OutTransitions[j].Description = fmt.Sprintf("transition from %s to %s", ctx.StateMachine.States[i].Code, ctx.StateMachine.States[i].OutTransitions[j].To)
			}
		}
	}
}

func (ctx *TokenContext) Valid() bool {
	v := ctx.Timeline.Valid()
	return v
}

func (ctx *TokenContext) IsActive() bool {
	return ctx.Timeline.IsInRange()
}

func (ctx *TokenContext) EvaluateInActions(tok *Token, params map[string]interface{}) ([]Action, error) {

	const semLogContext = "token-context::evaluate-in-actions"

	var state string
	if tok == nil || len(tok.Events) == 0 {
		state = StartEndState
	} else {
		state = tok.FindCurrentState()
	}

	log.Trace().Str("state", state).Msg(semLogContext + " relevant state")
	sd, err := ctx.StateMachine.FindStateDefinition(state)
	if err != nil {
		return nil, err
	}

	if len(sd.Actions) == 0 {
		log.Trace().Str("state", state).Msg(semLogContext + " no in actions found")
		return nil, nil
	}

	// Should check if token is in pending state?
	var eOpts []expression.Option
	if state != StartEndState {
		eOpts = []expression.Option{expression.WithVars(tok.Vars()), expression.WithMapInput(params)}
	} else {
		eOpts = []expression.Option{expression.WithMapInput(params)}
	}

	eCtx, err := expression.NewContext(eOpts...)
	if err != nil {
		return nil, err
	}

	acts, err := EvaluateActionDefinitions(sd.Actions, eCtx, ActionTypeIn, false)
	return acts, nil
}
