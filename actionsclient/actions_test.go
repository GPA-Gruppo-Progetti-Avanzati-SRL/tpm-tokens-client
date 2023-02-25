package actions_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/expression"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/hartracing"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/hartracing/filetracer"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
	actions "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-tokens-client/actionsclient"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestActions(t *testing.T) {

	closer, err := InitTracing(t)
	require.NoError(t, err)
	defer closer.Close()

	hc, err := InitHarTracing()
	require.NoError(t, err)
	defer hc.Close()

	cfg := []actions.Config{
		{
			Config: restclient.Config{
				RestTimeout: 0,
				SkipVerify:  false,
				Headers: []restclient.Header{
					{
						Name:  "request-id",
						Value: "{h:request-Id}",
					},
				},
				TraceGroupName:   "leas-cab-wfm-in-actions",
				TraceRequestName: "",
				RetryCount:       0,
				RetryWaitTime:    0,
				RetryMaxWaitTime: 0,
				RetryOnHttpError: nil,
				Span:             nil,
			},
			Id:   "id-of-action",
			Type: "bool",
			Host: actions.HostInfo{
				Scheme:   "http",
				HostName: "localhost",
				Port:     3001,
			},
			Method: http.MethodPost,
			Path:   "/api/v1/actions/id-of-action2",
		},
	}

	lks, err := actions.NewInstanceWithConfig(cfg)
	require.NoError(t, err)

	input := map[string]interface{}{
		"ssn": "MPRMLS62S21G337J",
	}
	exprCtx, err := expression.NewContext(
		expression.WithMapInput(input),
		expression.WithHeaders([]expression.NameValuePair{
			{
				Name: "request-id", Value: "my-req-id",
			},
		}))
	require.NoError(t, err)

	resp, err := lks.CallActions([]string{"id-of-action"}, exprCtx, input)
	require.NoError(t, err)

	t.Log(resp)
}

const (
	JAEGER_SERVICE_NAME = "JAEGER_SERVICE_NAME"
)

func InitHarTracing() (io.Closer, error) {
	trc, c, err := filetracer.NewTracer()
	if err != nil {
		return nil, err
	}

	hartracing.SetGlobalTracer(trc)
	return c, nil
}

func InitTracing(t *testing.T) (io.Closer, error) {

	if os.Getenv(JAEGER_SERVICE_NAME) == "" {
		t.Log("skipping jaeger config no vars in env.... (" + JAEGER_SERVICE_NAME + ")")
		return nil, nil
	}

	t.Log("initialize jaeger service " + os.Getenv(JAEGER_SERVICE_NAME))

	var tracer opentracing.Tracer
	var closer io.Closer

	jcfg, err := jaegercfg.FromEnv()
	if err != nil {
		log.Warn().Err(err).Msg("Unable to configure JAEGER from environment")
		return nil, err
	}

	tracer, closer, err = jcfg.NewTracer(
		jaegercfg.Logger(&jlogger{}),
		jaegercfg.Metrics(metrics.NullFactory),
	)
	if nil != err {
		log.Error().Err(err).Msg("Error in NewTracer")
		return nil, err
	}

	opentracing.SetGlobalTracer(tracer)

	return closer, nil
}

type jlogger struct{}

func (l *jlogger) Error(msg string) {
	log.Error().Msg("(jaeger) " + msg)
}

func (l *jlogger) Infof(msg string, args ...interface{}) {
	log.Info().Msgf("(jaeger) "+msg, args...)
}
