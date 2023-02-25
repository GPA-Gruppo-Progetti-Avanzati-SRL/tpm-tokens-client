package bridgeclient

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/har"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/hartracing"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
)

type ApiRequestContext struct {
	XAPIKey   string           `yaml:"x-api-key,omitempty" mapstructure:"x-api-key,omitempty" json:"x-api-key,omitempty"`
	RequestId string           `yaml:"request-id,omitempty" mapstructure:"request-id,omitempty" json:"request-id,omitempty"`
	Span      opentracing.Span `yaml:"-" mapstructure:"-" json:"-"`
	HarSpan   hartracing.Span  `yaml:"-" mapstructure:"-" json:"-"`
}

type APIRequestContextOption func(*ApiRequestContext)

func ApiRequestWithAutoRequestId() APIRequestContextOption {
	return func(ctx *ApiRequestContext) {
		ctx.RequestId = util.NewObjectId().String()
	}
}

func ApiRequestWithRequestId(reqId string) APIRequestContextOption {
	return func(ctx *ApiRequestContext) {
		if reqId == "" {
			reqId = util.NewObjectId().String()
			log.Warn().Msg("api-request: reqId set to empty string.... auto generated")
		}
		ctx.RequestId = reqId
	}
}

func ApiRequestWithApiKey(apiKey string) APIRequestContextOption {
	return func(ctx *ApiRequestContext) {
		if apiKey == "" {
			apiKey = util.NewObjectId().String()
			log.Warn().Msg("api-request: apiKey set to empty string.... auto generated")
		}
		ctx.XAPIKey = apiKey
	}
}

func ApiRequestWithSpan(span opentracing.Span) APIRequestContextOption {
	return func(ctx *ApiRequestContext) {
		ctx.Span = span
	}
}

func ApiRequestWithHarSpan(span hartracing.Span) APIRequestContextOption {
	return func(ctx *ApiRequestContext) {
		ctx.HarSpan = span
	}
}

func (arc *ApiRequestContext) getHeaders(ct string) []har.NameValuePair {

	const semLogContext = "bridge-client::get-headers"
	var nvp []har.NameValuePair

	if arc.RequestId == "" {
		arc.RequestId = util.NewObjectId().String()
	}
	nvp = append(nvp, har.NameValuePair{Name: RequestIdHeaderName, Value: arc.RequestId})

	if arc.XAPIKey != "" {
		nvp = append(nvp, har.NameValuePair{Name: ApiKeyHeaderName, Value: arc.XAPIKey})
	}

	if ct != "" {
		nvp = append(nvp, har.NameValuePair{Name: ContentTypeHeaderName, Value: ct})
	}

	return nvp
}

func NewApiRequestContext(opts ...APIRequestContextOption) ApiRequestContext {
	ar := ApiRequestContext{}
	for _, o := range opts {
		o(&ar)
	}

	return ar
}
