package campaignclient

import (
	"encoding/json"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-tokens-client/tokensclient"
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
	Code  string `yaml:"code,omitempty" mapstructure:"code,omitempty" json:"code,omitempty"`
	Ambit string `yaml:"ambit,omitempty" mapstructure:"ambit,omitempty" json:"ambit,omitempty"`
}

type Type struct {
	Code           string        `yaml:"code,omitempty" mapstructure:"code,omitempty" json:"code,omitempty"`
	Description    string        `yaml:"description,omitempty" mapstructure:"description,omitempty" json:"description,omitempty"`
	BannerTokenId  string        `yaml:"banner-token-id,omitempty" mapstructure:"banner-token-id,omitempty" json:"banner-token-id,omitempty"`
	Unique         bool          `yaml:"unique,omitempty" mapstructure:"unique,omitempty" json:"unique,omitempty"`
	PromoCode      string        `yaml:"promo,omitempty" mapstructure:"promo,omitempty" json:"promo,omitempty"`
	TargetProducts []ProductInfo `yaml:"target-products,omitempty" mapstructure:"target-products,omitempty" json:"target-products,omitempty"`
}

type AdditionalInfo struct {
	AltDescription   string `yaml:"alt-description,omitempty" mapstructure:"alt-description,omitempty" json:"alt-description,omitempty"`
	AwardDescription string `yaml:"award-description,omitempty" mapstructure:"award-description,omitempty" json:"award-description,omitempty"`
}

type LinkedResource struct {
	Type        string `yaml:"type,omitempty" mapstructure:"type,omitempty" json:"type,omitempty"`
	Name        string `yaml:"name,omitempty" mapstructure:"name,omitempty" json:"name,omitempty"`
	ContentType string `yaml:"content-type,omitempty" mapstructure:"content-type,omitempty" json:"content-type,omitempty"`
	Url         string `yaml:"url,omitempty" mapstructure:"url,omitempty" json:"url,omitempty"`
	Help        string `yaml:"help,omitempty" mapstructure:"help,omitempty" json:"help,omitempty"`
}

type Campaign struct {
	tokensclient.TokenContext `mapstructure:",squash"  yaml:",inline"`
	CampaignType              Type             `yaml:"type,omitempty" mapstructure:"type,omitempty" json:"type,omitempty"`
	Title                     string           `yaml:"title,omitempty" mapstructure:"title,omitempty" json:"title,omitempty"`
	Description               string           `yaml:"description,omitempty" mapstructure:"description,omitempty" json:"description,omitempty"`
	AddInfo                   AdditionalInfo   `yaml:"additional-info,omitempty" mapstructure:"additional-info,omitempty" json:"additional-info,omitempty"`
	Resources                 []LinkedResource `yaml:"resources,omitempty" mapstructure:"resources,omitempty" json:"resources,omitempty"`
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
	Id           string                `yaml:"id,omitempty" mapstructure:"id,omitempty" json:"id,omitempty"`
	Platform     string                `yaml:"platform,omitempty" mapstructure:"platform,omitempty" json:"platform,omitempty"`
	Version      string                `yaml:"version,omitempty" mapstructure:"version,omitempty" json:"version,omitempty"`
	Timeline     tokensclient.Timeline `yaml:"timeline,omitempty" mapstructure:"timeline,omitempty" json:"timeline,omitempty"`
	CampaignType Type                  `yaml:"type,omitempty" mapstructure:"type,omitempty" json:"type,omitempty"`
	Title        string                `yaml:"title,omitempty" mapstructure:"title,omitempty" json:"title,omitempty"`
	Description  string                `yaml:"description,omitempty" mapstructure:"description,omitempty" json:"description,omitempty"`
	AddInfo      AdditionalInfo        `yaml:"additional-info,omitempty" mapstructure:"additional-info,omitempty" json:"additional-info,omitempty"`
	Resources    []LinkedResource      `yaml:"resources,omitempty" mapstructure:"resources,omitempty" json:"resources,omitempty"`
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

func DeserializeContext(b []byte) (*Campaign, error) {
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
