package token

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"time"
)

type EventType string
type TokenType string

func (evtyp EventType) Scope() string {

	s := "undef"
	switch evtyp {
	case EventTypeCheck:
		s = string(EventTypeCheck)
	case EventTypeNext:
		s = string(EventTypeNext)
	case EventTypeCheckCreate:
		s = string(EventTypeCheck)
	case EventTypeCreate:
		s = string(EventTypeNext)
	default:
		log.Error().Interface("evt-type", evtyp).Msg("cannot determine operation scope from event type... reverting to next")
		s = string(EventTypeNext)
	}

	return s
}

const (
	EventTypeCheckCreate EventType = "check-create"
	EventTypeCreate      EventType = "create"
	EventTypeCheck       EventType = "check"
	EventTypeNext        EventType = "next"
	EventTypeCommit      EventType = "commit"
	EventTypeRollback    EventType = "rollback"
	EventTypeExpiration  EventType = "expired"

	TokenTypeStd    TokenType = "std"
	TokenTypeBanner TokenType = "banner"

	SysParamNameTokenId = "_tokId"

	semLogLabelTokenId = "token-id"
)

var NilTokenId = ""

/*
type TokenId string

func (t TokenId) String() string {
	return fmt.Sprintf("%s:%s:%s", t.CampaignId, t.TokenId, t.CheckDigit)
}
*/

// var TokenPatternRegexp = regexp.MustCompile("^([a-zA-Z]{6})\\:([a-zA-Z0-9]{1,16})\\:([a-zA-Z0-9])$")

/*
var TokenPatternRegExp = regexp.MustCompile("^([a-zA-Z]{6})([a-zA-Z0-9]{16})([a-zA-Z0-9])$")

func ParseTokenId(c string) (TokenId, bool) {

	matches := TokenPatternRegExp.FindAllSubmatch([]byte(c), -1)
	if len(matches) > 0 {
		t := TokenId{
			CampaignId: string(matches[0][1]),
			TokenId:    string(matches[0][2]),
			CheckDigit: string(matches[0][3]),
		}
		return t, true
	}

	return TokenId{}, false
}
*/

type TokenIdProvider interface {
	NewId(ctxId string, unique bool, action map[string]interface{}) (string, error)
}

type State struct {
	Code        string `yaml:"code,omitempty" mapstructure:"code,omitempty" json:"code,omitempty"`
	Description string `yaml:"description,omitempty" mapstructure:"description,omitempty" json:"description,omitempty"`
	Pending     bool   `yaml:"pending,omitempty" mapstructure:"pending,omitempty" json:"pending,omitempty"`
	LRAId       string `yaml:"lra-id,omitempty" mapstructure:"lra-id,omitempty" json:"lra-id,omitempty"`
}

type Timer struct {
	PKey            string           `yaml:"pkey,omitempty" mapstructure:"pkey,omitempty" json:"pkey,omitempty"`
	Id              string           `yaml:"id,omitempty" mapstructure:"id,omitempty" json:"id,omitempty"`
	Expires         string           `yaml:"expires,omitempty" mapstructure:"expires,omitempty" json:"expires,omitempty"`
	Outdated        bool             `yaml:"outdated,omitempty" mapstructure:"outdated,omitempty" json:"outdated,omitempty"`
	TimerDefinition *TimerDefinition `yaml:"definition,omitempty" mapstructure:"definition,omitempty" json:"definition,omitempty"`
}

type ProcessVars map[string]interface{}

type Event struct {
	RequestId      string      `yaml:"request-id,omitempty" mapstructure:"request-id,omitempty" json:"request-id,omitempty"`
	Name           string      `yaml:"name,omitempty" mapstructure:"name,omitempty" json:"name,omitempty"`
	Description    string      `yaml:"description,omitempty" mapstructure:"description,omitempty" json:"description,omitempty"`
	Typ            EventType   `yaml:"type,omitempty" mapstructure:"type,omitempty" json:"type,omitempty"`
	State          State       `yaml:"state,omitempty" mapstructure:"state,omitempty" json:"state,omitempty"`
	Ts             string      `yaml:"ts,omitempty" mapstructure:"ts,omitempty" json:"ts,omitempty"`
	ExpiryTs       string      `yaml:"expiry-ts,omitempty" mapstructure:"expiry-ts,omitempty" json:"expiry-ts,omitempty"`
	Vars           ProcessVars `yaml:"vars,omitempty" mapstructure:"vars,omitempty" json:"vars,omitempty"`
	Actions        []Action    `yaml:"actions,omitempty" mapstructure:"actions,omitempty" json:"actions,omitempty"`
	Bearers        []BearerRef `yaml:"bearers,omitempty" mapstructure:"bearers,omitempty" json:"bearers,omitempty"`
	TimerReference *Timer      `yaml:"timer-ref,omitempty" mapstructure:"timer-ref,omitempty" json:"timer-ref,omitempty"`
}

func (evt *Event) FindAction(actionId string, actionType ActionType) (Action, bool) {
	for _, a := range evt.Actions {
		if a.ActionType == actionType && a.ActionId == actionId {
			return a, true
		}
	}
	return Action{}, false
}

func (evt *Event) IsPending() bool {
	return evt.State.Pending
}

type Token struct {
	Pkey     string                 `yaml:"pkey,omitempty" mapstructure:"pkey" json:"pkey,omitempty"`
	Id       string                 `yaml:"id,omitempty" mapstructure:"id,omitempty" json:"id,omitempty"`
	Typ      TokenType              `yaml:"type,omitempty" mapstructure:"type,omitempty" json:"type,omitempty"`
	CtxId    string                 `yaml:"ctx-id,omitempty" mapstructure:"ctx-id,omitempty" json:"ctx-id,omitempty"`
	Events   []Event                `yaml:"events,omitempty" mapstructure:"events,omitempty" json:"events,omitempty"`
	Metadata map[string]interface{} `yaml:"metadata,omitempty" mapstructure:"metadata,omitempty" json:"metadata,omitempty"`
	TTL      int                    `yaml:"ttl,omitempty" mapstructure:"ttl,omitempty" json:"ttl,omitempty"`
}

func DeserializeToken(b []byte) (*Token, error) {
	tok := Token{}
	err := json.Unmarshal(b, &tok)
	if err != nil {
		return nil, err
	}

	return &tok, nil
}

func (tok *Token) TokenId() (string, error) {
	return tok.Id, nil

	/*
		tid, ok := ParseTokenId(tok.Id)
		if !ok {
			return tid, fmt.Errorf("cannot parse %s", tok.Id)
		}

		return tid, nil
	*/
}

func (tok *Token) ToJSON() ([]byte, error) {
	return json.Marshal(tok)
}

func (tok *Token) MustToJSON() []byte {
	b, err := json.Marshal(tok)
	if err != nil {
		panic(err)
	}

	return b
}

func (tok *Token) FindOutdatedTimers() []*Timer {
	var tms []*Timer
	for i := 0; i < len(tok.Events)-1; i++ {
		if tok.Events[i].TimerReference != nil {
			tms = append(tms, tok.Events[i].TimerReference)
		}
	}

	return tms
}

func (tok *Token) MarkTimersAsOutdated() {
	for i := 0; i < len(tok.Events); i++ {
		if tok.Events[i].TimerReference != nil {
			tok.Events[i].TimerReference.Outdated = true
			tok.Events[i].TimerReference.TimerDefinition = nil
		}
	}
}

func (tok *Token) IsPending() bool {
	ndx := tok.FindLastEventIndex()
	if ndx >= 0 {
		return tok.Events[ndx].IsPending()
	}

	return false
}

func (tok *Token) FindEventIndexByState(st string) int {

	foundNdx := -1
	for i, e := range tok.Events {
		if e.State.Code == st {
			foundNdx = i
		}
	}

	return foundNdx
}

func (tok *Token) FindLastEventIndex() int {
	return len(tok.Events) - 1
}

func (tok *Token) FindCurrentState() string {
	if len(tok.Events) > 0 {
		return tok.Events[len(tok.Events)-1].State.Code
	}

	log.Warn().Str(semLogLabelTokenId, tok.Id).Msg("token unknown state")
	return ""
}

func (tok *Token) FindEventIndexByRequestId(reqId string) int {

	for i := len(tok.Events) - 1; i >= 0; i-- {
		if tok.Events[i].RequestId == reqId {
			return i
		}
	}

	return -1
}

func (tok *Token) IsExpired(timelineMode string) bool {
	lastEvt := tok.FindLastEventIndex()
	if lastEvt < 0 {
		log.Warn().Str(semLogLabelTokenId, tok.Id).Msg("token is empty")
		return true
	}

	expTs := tok.Events[lastEvt].ExpiryTs
	if expTs == "" {
		return false
	}

	if timelineMode == ExpirationModeTimestamp {
		expTm, err := time.Parse(time.RFC3339, expTs)
		if err != nil {
			log.Error().Err(err).Str(semLogLabelTokenId, tok.Id).Str("expiry-ts", expTs).Msg("invalid token expiry ts")
			return true
		}

		return time.Now().After(expTm)
	}

	return expTs < time.Now().Format("20060102")
}

func (tok *Token) Vars() ProcessVars {
	if len(tok.Events) == 0 {
		return nil
	}
	lastEvt := tok.Events[len(tok.Events)-1]
	return lastEvt.Vars
}
