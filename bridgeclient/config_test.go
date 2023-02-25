package bridgeclient_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/hartracing"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/hartracing/filetracer"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-tokens-client/bridgeclient"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
	"io"
	"os"
	"testing"
)

var cliConfig = bridgeclient.Config{
	Config: restclient.Config{
		RestTimeout:      0,
		SkipVerify:       false,
		Headers:          nil,
		TraceGroupName:   "bridge-client-api",
		TraceRequestName: "",
		RetryCount:       0,
		RetryWaitTime:    0,
		RetryMaxWaitTime: 0,
		RetryOnHttpError: nil,
	},
	Host: bridgeclient.HostInfo{
		Scheme:   "http",
		HostName: "localhost",
		Port:     8085,
	},
	Endpoints: []bridgeclient.EndpointDefinition{
		{
			Id:  "new-id",
			Url: "/api/v1/campaigns/{context-id}",
		},
	},
}

func TestMain(m *testing.M) {
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestBridgeClient(t *testing.T) {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	c, err := InitTracing(t)
	require.NoError(t, err)
	if c != nil {
		defer c.Close()
	}

	cHar, err := InitHarTracing(t)
	require.NoError(t, err)
	if cHar != nil {
		defer cHar.Close()
	}

	cli, err := bridgeclient.NewClient(&cliConfig)
	require.NoError(t, err)
	defer cli.Close()

	executeTestBridgeClient(t, cli)
}

func executeTestBridgeClient(t *testing.T, cli *bridgeclient.Client) {
	apiRequestCtx := bridgeclient.NewApiRequestContext(bridgeclient.ApiRequestWithApiKey("ApiKeyLeasCabBridge"))

	resp, err := cli.NewId(apiRequestCtx, "BPMIFI", true, map[string]interface{}{"cf": "MPRMLS62S21G337J"})
	require.NoError(t, err)
	t.Log(resp)
}

const (
	JAEGER_SERVICE_NAME = "JAEGER_SERVICE_NAME"
)

func InitHarTracing(t *testing.T) (io.Closer, error) {
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
