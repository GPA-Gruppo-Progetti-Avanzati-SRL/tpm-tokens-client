package facts

import "encoding/json"

type Fact struct {
	Id         string                 `yaml:"id,omitempty" mapstructure:"id,omitempty" json:"id,omitempty"`
	Group      string                 `yaml:"pkey,omitempty" mapstructure:"pkey,omitempty" json:"pkey,omitempty"`
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
