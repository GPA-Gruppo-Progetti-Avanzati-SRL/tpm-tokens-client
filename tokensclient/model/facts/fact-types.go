package facts

import "encoding/json"

type Fact struct {
	Class      string                 `yaml:"class,omitempty" mapstructure:"class,omitempty" json:"class,omitempty"`
	Group      string                 `yaml:"group,omitempty" mapstructure:"group,omitempty" json:"group,omitempty"`
	Id         string                 `yaml:"id,omitempty" mapstructure:"id,omitempty" json:"id,omitempty"`
	PKey       string                 `yaml:"pkey,omitempty" mapstructure:"pkey,omitempty" json:"pkey,omitempty"`
	CtxId      string                 `yaml:"ctx-id,omitempty" mapstructure:"ctx-id,omitempty" json:"ctx-id,omitempty"`
	TokenId    string                 `yaml:"token-id,omitempty" mapstructure:"token-id,omitempty" json:"token-id,omitempty"`
	Properties map[string]interface{} `yaml:"properties,omitempty" mapstructure:"properties,omitempty" json:"properties,omitempty"`
	TTL        int                    `yaml:"ttl,omitempty" mapstructure:"ttl,omitempty" json:"ttl,omitempty"`
}

func (ctx *Fact) ToJSON() ([]byte, error) {
	return json.Marshal(ctx)
}

func (ctx *Fact) MustToJSON() []byte {
	b, err := json.Marshal(ctx)
	if err != nil {
		panic(err)
	}

	return b
}

func DeserializeFact(b []byte) (*Fact, error) {
	ctx := Fact{}
	err := json.Unmarshal(b, &ctx)
	if err != nil {
		return nil, err
	}

	return &ctx, nil
}

// FactsQueryResponse (vedi tpm-tokens/tokensazstore/fact.go)
type FactsQueryResponse struct {
	RespRid   string `json:"_rid" yaml:"_rid"`
	RespCount int    `json:"_count" yaml:"_count"`
	Documents []Fact `json:"documents,omitempty" yaml:"documents,omitempty"`
}

func DeserializeFactsQueryResponse(b []byte) (*FactsQueryResponse, error) {
	ctx := FactsQueryResponse{}
	err := json.Unmarshal(b, &ctx)
	if err != nil {
		return nil, err
	}

	return &ctx, nil
}
