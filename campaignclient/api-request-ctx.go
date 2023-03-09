package campaignclient

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
	LRAId     string           `yaml:"lra-id,omitempty" mapstructure:"lra-id,omitempty" json:"lra-id,omitempty"`
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
	const semLogContext = "tokens-client::request-req-id"
	return func(ctx *ApiRequestContext) {
		if reqId == "" {
			reqId = util.NewObjectId().String()
			log.Info().Msg(semLogContext + " reqId set to empty string.... auto generated")
		}
		ctx.RequestId = reqId
	}
}

func ApiRequestWithApiKey(apiKey string) APIRequestContextOption {
	const semLogContext = "tokens-client::request-with-api-key"
	return func(ctx *ApiRequestContext) {
		if apiKey == "" {
			apiKey = util.NewObjectId().String()
			log.Info().Msg(semLogContext + " apiKey set to empty string.... auto generated")
		}
		ctx.XAPIKey = apiKey
	}
}

func ApiRequestWithAutoLraId() APIRequestContextOption {
	return func(ctx *ApiRequestContext) {
		ctx.LRAId = util.NewObjectId().String()
	}
}

func ApiRequestWithLraId(lraId string) APIRequestContextOption {
	const semLogContext = "tokens-client::request-with-lra-id"
	return func(ctx *ApiRequestContext) {
		if lraId == "" {
			log.Info().Msg(semLogContext + " lraId set to empty string.... no lraId set")
		}
		ctx.LRAId = lraId
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

	const semLogContext = "tpm-tokens-client::get-headers"
	var nvp []har.NameValuePair

	if arc.RequestId == "" {
		arc.RequestId = util.NewObjectId().String()
	}
	nvp = append(nvp, har.NameValuePair{Name: RequestIdHeaderName, Value: arc.RequestId})

	if arc.LRAId != "" {
		nvp = append(nvp, har.NameValuePair{Name: LraHttpContextHeaderName, Value: arc.LRAId})
	}

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
