package campaignclient

import (
	"encoding/json"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-tokens-client/tokensclient/model/token"
	"strings"
	"time"
)

const (
	CosPartitionKey = "campaign"
)

type Token struct {
	CampaignId string
	TokenId    string
	CheckDigit string
}

func (t *Token) String() string {
	return fmt.Sprintf("%s:%s:%s", t.CampaignId, t.TokenId, t.CheckDigit)
}

type ProductInfo struct {
	Code string `yaml:"code,omitempty" mapstructure:"code,omitempty" json:"code,omitempty"`
}

type TokenMode string

const (
	TokenModeBatch  TokenMode = "batch"
	TokenModeOnline TokenMode = "online"
)

type Type struct {
	Code           string        `yaml:"code,omitempty" mapstructure:"code,omitempty" json:"code,omitempty"`
	Description    string        `yaml:"description,omitempty" mapstructure:"description,omitempty" json:"description,omitempty"`
	Unique         bool          `yaml:"unique,omitempty" mapstructure:"unique,omitempty" json:"unique,omitempty"`
	CpqServiceCode string        `yaml:"cpq-service-code,omitempty" mapstructure:"cpq-service-code,omitempty" json:"cpq-service-code,omitempty"`
	TokenMode      TokenMode     `yaml:"token-mode,omitempty" mapstructure:"token-mode,omitempty" json:"token-mode,omitempty"`
	TargetProducts []ProductInfo `yaml:"target-products,omitempty" mapstructure:"target-products,omitempty" json:"target-products,omitempty"`
}

type AdditionalInfo struct {
	AltDescription   string `yaml:"alt-description,omitempty" mapstructure:"alt-description,omitempty" json:"alt-description,omitempty"`
	AwardDescription string `yaml:"award-description,omitempty" mapstructure:"award-description,omitempty" json:"award-description,omitempty"`
}

type LinkedResourceLocation struct {
	Type string `yaml:"type,omitempty" mapstructure:"type,omitempty" json:"type,omitempty"`
	Url  string `yaml:"url,omitempty" mapstructure:"url,omitempty" json:"url,omitempty"`
}

type LinkedResource struct {
	Type        string                   `yaml:"type,omitempty" mapstructure:"type,omitempty" json:"type,omitempty"`
	Name        string                   `yaml:"name,omitempty" mapstructure:"name,omitempty" json:"name,omitempty"`
	ContentType string                   `yaml:"content-type,omitempty" mapstructure:"content-type,omitempty" json:"content-type,omitempty"`
	Locations   []LinkedResourceLocation `yaml:"locations,omitempty" mapstructure:"locations,omitempty" json:"locations,omitempty"`
	Help        string                   `yaml:"help,omitempty" mapstructure:"help,omitempty" json:"help,omitempty"`
	Properties  map[string]interface{}   `yaml:"properties,omitempty" mapstructure:"properties,omitempty" json:"properties,omitempty"`
}

type Filters struct {
	Canale   string `yaml:"canale,omitempty" mapstructure:"canale,omitempty" json:"canale,omitempty"`
	Servizio string `yaml:"servizio,omitempty" mapstructure:"servizio,omitempty" json:"servizio,omitempty"`
	Prodotto string `yaml:"prodotto,omitempty" mapstructure:"prodotto,omitempty" json:"prodotto,omitempty"`
	Fase     string `yaml:"fase,omitempty" mapstructure:"fase,omitempty" json:"fase,omitempty"`
	Timing   string `yaml:"-" mapstructure:"-" json:"-"`
}

type Campaign struct {
	token.TokenContext `mapstructure:",squash"  yaml:",inline"`
	Filters            Filters          `yaml:"filters,omitempty" mapstructure:"filters,omitempty" json:"filters,omitempty"`
	CampaignType       Type             `yaml:"type,omitempty" mapstructure:"type,omitempty" json:"type,omitempty"`
	Title              string           `yaml:"title,omitempty" mapstructure:"title,omitempty" json:"title,omitempty"`
	Description        string           `yaml:"description,omitempty" mapstructure:"description,omitempty" json:"description,omitempty"`
	AddInfo            AdditionalInfo   `yaml:"additional-info,omitempty" mapstructure:"additional-info,omitempty" json:"additional-info,omitempty"`
	Resources          []LinkedResource `yaml:"resources,omitempty" mapstructure:"resources,omitempty" json:"resources,omitempty"`
}

func (c *Campaign) Info() CampaignInfo {
	return CampaignInfo{
		Id:           c.Id,
		Platform:     c.Platform,
		Version:      c.Version,
		Timeline:     c.Timeline,
		CampaignType: c.CampaignType,
		Title:        c.Title,
		Description:  c.Description,
		AddInfo:      c.AddInfo,
		Resources:    c.Resources,
	}
}

type CampaignInfo struct {
	Id           string           `yaml:"id,omitempty" mapstructure:"id,omitempty" json:"id,omitempty"`
	Platform     string           `yaml:"platform,omitempty" mapstructure:"platform,omitempty" json:"platform,omitempty"`
	Version      string           `yaml:"version,omitempty" mapstructure:"version,omitempty" json:"version,omitempty"`
	Timeline     token.Timeline   `yaml:"timeline,omitempty" mapstructure:"timeline,omitempty" json:"timeline,omitempty"`
	Filters      Filters          `yaml:"filters,omitempty" mapstructure:"filters,omitempty" json:"filters,omitempty"`
	CampaignType Type             `yaml:"type,omitempty" mapstructure:"type,omitempty" json:"type,omitempty"`
	Title        string           `yaml:"title,omitempty" mapstructure:"title,omitempty" json:"title,omitempty"`
	Description  string           `yaml:"description,omitempty" mapstructure:"description,omitempty" json:"description,omitempty"`
	AddInfo      AdditionalInfo   `yaml:"additional-info,omitempty" mapstructure:"additional-info,omitempty" json:"additional-info,omitempty"`
	Resources    []LinkedResource `yaml:"resources,omitempty" mapstructure:"resources,omitempty" json:"resources,omitempty"`
}

func (c *CampaignInfo) Accept(criteria *Filters) bool {

	rc := true
	if rc && !checkTiming(c.Timeline, criteria.Timing) {
		rc = false
	}

	if rc && !checkFilter(c.Filters.Canale, criteria.Canale) {
		rc = false
	}

	if rc && !checkFilter(c.Filters.Servizio, criteria.Servizio) {
		rc = false
	}

	if rc && !checkFilter(c.Filters.Prodotto, criteria.Prodotto) {
		rc = false
	}

	if rc && !checkFilter(c.Filters.Fase, criteria.Fase) {
		rc = false
	}

	return rc
}

func checkFilter(val string, criteria string) bool {
	if val == "*" || criteria == "" {
		return true
	}

	val = strings.ToLower(val)

	if strings.Index(val, ",") >= 0 {
		arr := strings.Split(val, ",")
		for _, s := range arr {
			if s == criteria {
				return true
			}
		}

		return false
	}

	return strings.ToLower(val) == criteria
}

func checkTiming(val token.Timeline, criteria string) bool {
	if criteria == "" {
		return true
	}

	today := time.Now().Format("20060102")
	rc := false
	switch criteria {
	case "next":
		if today < val.StartDate {
			rc = true
		}
	case "current":
		if today >= val.StartDate && today <= val.EndDate {
			rc = true
		}
	case "past":
		if today > val.EndDate {
			rc = true
		}
	}

	return rc
}

/*
func (c *Campaign) SetTokenContext(ctx *tokensclient.TokenContext) {
	c.Id = ctx.Id
	c.Pkey = ctx.Pkey
	c.Platform = ctx.Platform
	c.Version = ctx.Version
	c.Timeline = ctx.Timeline
	c.Suspended = ctx.Suspended
	c.StateMachine = ctx.StateMachine
	c.TokenIdProviderType = ctx.TokenIdProviderType
}
*/

func (c *Campaign) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Campaign) MustToJSON() []byte {
	b, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}

	return b
}

func WellFormCampaignId(id string) string {
	return strings.ToUpper(id)
}

func DeserializeCampaign(b []byte) (*Campaign, error) {
	ctx := Campaign{}
	err := json.Unmarshal(b, &ctx)
	if err != nil {
		return nil, err
	}

	ctx.Pkey = CosPartitionKey
	return &ctx, nil
}

func (c *Campaign) Valid() bool {
	v := c.TokenContext.Valid()
	return v
}

func (c *Campaign) IsActive() bool {
	today := time.Now().Format("20060102")
	if today >= c.Timeline.StartDate && today <= c.Timeline.EndDate {
		return true
	}

	return false
}
