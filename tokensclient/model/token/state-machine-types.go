package token

import (
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/expression"
)

type StateType string

const (
	StateTransient StateType = "transient"
	StateFinal     StateType = "final"
	StateExpired   StateType = "expired"
	StateStd       StateType = "std"

	StartEndState = "[*]"

	ExpiryTsProcessVariable = "expiry-ts"
)

type ActionType string

const (
	ActionTypeNewId = "new-id"
	ActionTypeIn    = "action-in"
	ActionTypeOut   = "action-out"
)

type ActionDefinition struct {
	ActionId   string                 `yaml:"id,omitempty" mapstructure:"id,omitempty" json:"id,omitempty"`
	ActionType ActionType             `yaml:"type,omitempty" mapstructure:"type,omitempty" json:"type,omitempty"`
	Properties map[string]interface{} `yaml:"properties,omitempty" mapstructure:"properties,omitempty" json:"properties,omitempty"`
}

type Action ActionDefinition

type TTLDefinition struct {
	Value string `yaml:"value,omitempty" mapstructure:"value,omitempty" json:"value,omitempty"`
}

type CodeDescriptionPair struct {
	Code        string `yaml:"code,omitempty" mapstructure:"code,omitempty" json:"code,omitempty"`
	Description string `yaml:"description,omitempty" mapstructure:"description,omitempty" json:"description,omitempty"`
}

type ProcessVarDefinition struct {
	Name        string `yaml:"name,omitempty" mapstructure:"name,omitempty" json:"name,omitempty"`
	Description string `yaml:"description,omitempty" mapstructure:"description,omitempty" json:"description,omitempty"`
	Value       string `yaml:"value,omitempty" mapstructure:"value,omitempty" json:"value,omitempty"`
}

type StateDefinition struct {
	Code        string `yaml:"code,omitempty" mapstructure:"code,omitempty" json:"code,omitempty"`
	Description string `yaml:"description,omitempty" mapstructure:"description,omitempty" json:"description,omitempty"`
	// Help                     string              `yaml:"help,omitempty" mapstructure:"help,omitempty" json:"help,omitempty"`
	StateType      StateType           `yaml:"type,omitempty" mapstructure:"type,omitempty" json:"type,omitempty"`
	OutTransitions []Transition        `yaml:"transitions,omitempty" mapstructure:"transitions,omitempty" json:"transitions,omitempty"`
	Actions        []ActionDefinition  `yaml:"in-actions,omitempty" mapstructure:"in-actions,omitempty" json:"in-actions,omitempty"`
	Help           CodeDescriptionPair `yaml:"help,omitempty" mapstructure:"help,omitempty" json:"help,omitempty"`
}

type Rule struct {
	Expression string              `yaml:"expr,omitempty" mapstructure:"expr,omitempty" json:"expr,omitempty"`
	Help       CodeDescriptionPair `yaml:"help,omitempty" mapstructure:"help,omitempty" json:"help,omitempty"`
}

type Property struct {
	Name           string              `yaml:"name,omitempty" mapstructure:"name,omitempty" json:"name,omitempty"`
	ValidationRule string              `yaml:"validation-rule,omitempty" mapstructure:"validation-rule,omitempty" json:"validation-rule,omitempty"`
	Help           CodeDescriptionPair `yaml:"help,omitempty" mapstructure:"help,omitempty" json:"help,omitempty"`
	Scope          string              `yaml:"scope,omitempty" mapstructure:"scope,omitempty" json:"scope,omitempty"`
}

type BearerRef struct {
	Id   string `yaml:"id,omitempty" mapstructure:"id,omitempty" json:"id,omitempty"`
	Role string `yaml:"role,omitempty" mapstructure:"role,omitempty" json:"role,omitempty"`
}

type Transition struct {
	Name                  string                 `yaml:"name,omitempty" mapstructure:"name,omitempty" json:"name,omitempty"`
	To                    string                 `yaml:"to,omitempty" mapstructure:"to,omitempty" json:"to,omitempty"`
	Order                 int                    `yaml:"order,omitempty" mapstructure:"order,omitempty" json:"order,omitempty"`
	Properties            []Property             `yaml:"properties,omitempty" mapstructure:"properties,omitempty" json:"properties,omitempty"`
	Rules                 []Rule                 `yaml:"rules,omitempty" mapstructure:"rules,omitempty" json:"rules,omitempty"`
	ProcessVarDefinitions []ProcessVarDefinition `yaml:"process-vars,omitempty" mapstructure:"process-vars,omitempty" json:"process-vars,omitempty"`
	TTL                   TTLDefinition          `yaml:"ttl,omitempty" mapstructure:"ttl,omitempty" json:"ttl,omitempty"`
	Actions               []ActionDefinition     `yaml:"out-actions,omitempty" mapstructure:"out-actions,omitempty" json:"out-actions,omitempty"`
	Bearers               []BearerRef            `yaml:"bearers,omitempty" mapstructure:"bearers,omitempty" json:"bearers,omitempty"`
	Description           string                 `yaml:"description,omitempty" mapstructure:"description,omitempty" json:"description,omitempty"`
	TimerDefinition       *TimerDefinition       `yaml:"timer-def,omitempty" mapstructure:"timer-def,omitempty" json:"timer-def,omitempty"`
}

type Diagram struct {
	ContentType string `yaml:"content-type,omitempty" mapstructure:"content-type,omitempty" json:"content-type,omitempty"`
	Data        string `yaml:"data,omitempty" mapstructure:"data,omitempty" json:"data,omitempty"`
}

type StateMachine struct {
	Diagram          *Diagram          `yaml:"diagram,omitempty" mapstructure:"diagram,omitempty" json:"diagram,omitempty"`
	States           []StateDefinition `yaml:"states,omitempty" mapstructure:"states,omitempty" json:"states,omitempty"`
	CatchTransitions []Transition      `yaml:"catch-transitions,omitempty" mapstructure:"catch-transitions,omitempty" json:"catch-transitions,omitempty"`
}

func (sm *StateMachine) FindStateDefinition(code string) (StateDefinition, error) {
	for _, s := range sm.States {
		if s.Code == code {
			return s, nil
		}
	}

	return StateDefinition{}, NewTokError(TokenErrorContextDefinition, fmt.Sprintf("cannot find definition of state: %s", code))
}

func EvaluateActionDefinitions(actions []ActionDefinition, eCtx *expression.Context, actionType ActionType, takeFirstOnly bool) ([]Action, error) {
	if len(actions) == 0 {
		return nil, nil
	}

	var acts []Action
	for _, a := range actions {

		if a.ActionType == actionType {
			v, err := EvalProperties(eCtx, a.Properties)
			if err != nil {
				return nil, err
			}

			acts = append(acts, Action{ActionId: a.ActionId, Properties: v, ActionType: a.ActionType})
			if takeFirstOnly {
				return acts, nil
			}
		}
	}

	return acts, nil
}

func EvalProperties(eCtx *expression.Context, propsDefinition map[string]interface{}) (map[string]interface{}, error) {

	var err error
	if len(propsDefinition) > 0 {
		pv := make(map[string]interface{})
		for n, v := range propsDefinition {
			s := fmt.Sprintf("%v", v)
			if _, ok := v.(string); ok {
				s = v.(string)
			}
			pv[n], err = eCtx.EvalOne(s)

			if err != nil {
				return nil, NewTokError(TokenErrorExpressionEvaluation, err.Error())
			}
		}

		return pv, nil
	}

	return nil, nil
}
