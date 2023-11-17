package bearer

import (
	"encoding/json"
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

func NewBearer(bearerId, contextId string) Bearer {
	return Bearer{Id: Id(bearerId, contextId), Pkey: bearerId, TokenContextId: contextId, TTL: -1}
}

func Id(bearerId, contextId string) string {
	return strings.Join([]string{bearerId, contextId}, "-")
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
