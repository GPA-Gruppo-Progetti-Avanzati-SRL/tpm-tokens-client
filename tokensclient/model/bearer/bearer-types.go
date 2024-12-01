package bearer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"

	"strings"
)

const (
	MailPropertyName      = "mail"
	FirstNamePropertyName = "first-name"
	LastNamePropertyName  = "last-name"
)

type TokenRef struct {
	Id   string `yaml:"id,omitempty" mapstructure:"id,omitempty" json:"id,omitempty"`
	Role string `yaml:"role,omitempty" mapstructure:"role,omitempty" json:"role,omitempty"`
}

type Bearer struct {
	Id             string                 `yaml:"id,omitempty" mapstructure:"id,omitempty" json:"id,omitempty"`
	Pkey           string                 `yaml:"pkey,omitempty" mapstructure:"pkey" json:"pkey,omitempty"`
	TokenContextId string                 `yaml:"tok-ctx-id,omitempty" mapstructure:"tok-ctx-id,omitempty" json:"tok-ctx-id,omitempty"`
	ActorId        string                 `yaml:"actor-id,omitempty" mapstructure:"actor-id,omitempty" json:"actor-id,omitempty"`
	ActorScope     string                 `yaml:"actor-scope,omitempty" mapstructure:"actor-scope,omitempty" json:"actor-scope,omitempty"`
	Origin         string                 `yaml:"origin,omitempty" mapstructure:"origin,omitempty" json:"origin,omitempty"`
	TokenRefs      []TokenRef             `yaml:"tok-refs,omitempty" mapstructure:"tok-refs,omitempty" json:"tok-refs,omitempty"`
	Properties     map[string]interface{} `yaml:"properties,omitempty" mapstructure:"properties,omitempty" json:"properties,omitempty"`
	TTL            int                    `yaml:"ttl,omitempty" mapstructure:"ttl,omitempty" json:"ttl,omitempty"`
}

func (ber *Bearer) ToJSON() ([]byte, error) {
	return json.Marshal(ber)
}

func (ber *Bearer) MustToJSON() []byte {
	b, err := json.Marshal(ber)
	if err != nil {
		panic(err)
	}

	return b
}

func DeserializeBearer(b []byte) (*Bearer, error) {
	ctx := Bearer{}
	err := json.Unmarshal(b, &ctx)
	if err != nil {
		return nil, err
	}

	return &ctx, nil
}

const (
	ActorScopeMatrixParamValue = ";scope="
)

func NewBearer(actorId, actorScope, contextId string) Bearer {
	var b Bearer
	if strings.Index(actorId, ActorScopeMatrixParamValue) >= 0 {
		actorId, actorScope, _, _ = ParseBearerId(actorId)
		b = Bearer{Id: Id(actorId, actorScope, contextId), Pkey: actorId, ActorId: actorScope, TokenContextId: contextId, TTL: -1}
	} else {
		b = Bearer{Id: Id(actorId, actorScope, contextId), Pkey: actorId, ActorId: actorId, ActorScope: actorScope, TokenContextId: contextId, TTL: -1}
	}
	return b
}

func ActorIdWithScope(actorId, actorScope string) string {
	if strings.Index(actorId, ActorScopeMatrixParamValue) < 0 && actorScope != "" {
		actorId = fmt.Sprintf("%s%s%s", actorId, ActorScopeMatrixParamValue, actorScope)
	}

	return actorId
}

func Id(actorId, actorScope, contextId string) string {
	if strings.Index(actorId, ActorScopeMatrixParamValue) < 0 && actorScope != "" {
		actorId = fmt.Sprintf("%s%s%s", actorId, ActorScopeMatrixParamValue, actorScope)
	}
	return strings.Join([]string{actorId, contextId}, "-")
}

func ParseBearerId(bearerId string) (string, string, string, error) {
	const semLogContext = "bearer::parse-bearer-id"

	comps := strings.Split(bearerId, "-")
	if len(comps) != 2 {
		err := errors.New("invalid bearer id")
		log.Error().Err(err).Str("bearer-id", bearerId).Msg(semLogContext)
		return bearerId, "", "", err
	}

	actorId := comps[0]
	campaignId := comps[1]
	var actorScope string
	if ndx := strings.Index(comps[0], ActorScopeMatrixParamValue); ndx >= 0 {
		actorScope = comps[0][ndx+len(ActorScopeMatrixParamValue):]
		actorId = comps[0][:ndx]
	}
	return actorId, actorScope, campaignId, nil
}

func WellFormBearerId(id string) string {
	return strings.ToUpper(id)
}

func (ber *Bearer) HasToken(tokId string) bool {
	for _, t := range ber.TokenRefs {
		if t.Id == tokId {
			return true
		}
	}

	return false
}

func (ber *Bearer) AddToken(tokId string, role string) bool {

	const semLogContext = "bearer::add-token"
	for _, t := range ber.TokenRefs {
		if t.Id == tokId {
			return false
		}
	}

	if role == "" {
		role = "primary"
	}

	if role != "primary" && role != "secondary" {
		log.Error().Str("role", role).Str("token-id", tokId).Msg(semLogContext + " unsupported role")
	}

	ber.TokenRefs = append(ber.TokenRefs, TokenRef{Id: tokId, Role: role})
	return true
}

func (ber *Bearer) RemoveToken(tokId string) bool {

	foundNdx := 1
	for i, t := range ber.TokenRefs {
		if t.Id == tokId {
			foundNdx = i
			break
		}
	}

	if foundNdx >= 0 {
		if len(ber.TokenRefs) == 1 {
			ber.TokenRefs = nil
		} else {
			if foundNdx == len(ber.TokenRefs)-1 {
				ber.TokenRefs = ber.TokenRefs[0:foundNdx]
			} else {
				if foundNdx == 0 {
					ber.TokenRefs = ber.TokenRefs[1:]
				} else {
					ber.TokenRefs = append(ber.TokenRefs[0:foundNdx], ber.TokenRefs[foundNdx+1:]...)
				}
			}
		}

		return true
	}

	return false
}

type BearersQueryResponse struct {
	RespRid   string   `json:"_rid" yaml:"_rid"`
	RespCount int      `json:"_count" yaml:"_count"`
	Documents []Bearer `json:"documents,omitempty" yaml:"documents,omitempty"`
}

func DeserializeBearersQueryResponse(b []byte) (*BearersQueryResponse, error) {
	ctx := BearersQueryResponse{}
	err := json.Unmarshal(b, &ctx)
	if err != nil {
		return nil, err
	}

	return &ctx, nil
}
