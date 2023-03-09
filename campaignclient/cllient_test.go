package campaignclient_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/hartracing"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/hartracing/filetracer"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-tokens-client/campaignclient"

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

var cliConfig = campaignclient.Config{
	Config: restclient.Config{
		RestTimeout:      0,
		SkipVerify:       false,
		Headers:          nil,
		TraceGroupName:   "campaign-client-api",
		TraceRequestName: "",
		RetryCount:       0,
		RetryWaitTime:    0,
		RetryMaxWaitTime: 0,
		RetryOnHttpError: nil,
	},
	Host: campaignclient.HostInfo{
		Scheme:   "http",
		HostName: "localhost",
		Port:     8082,
	},
}

func TestMain(m *testing.M) {
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestCampaignClient(t *testing.T) {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	c, err := InitTracing(t)
	require.NoError(t, err)
	if c != nil {
		defer c.Close()
	}

	cli, err := campaignclient.NewTokensApiClient(&cliConfig)
	require.NoError(t, err)
	defer cli.Close()

	executeTestCampaignClient(t, cli, &campaignTestCase001)
}

func executeTestCampaignClient(t *testing.T, cli *campaignclient.Client, campaignTestCase *campaignclient.Campaign) {
	apiRequestCtx := campaignclient.NewApiRequestContext(campaignclient.ApiRequestWithApiKey("ApiKeyWfm"))

	campaignTestCase.Id = "TESTCTX"
	resp, err := cli.NewCampaign(apiRequestCtx, campaignTestCase, "")
	require.NoError(t, err)
	t.Log(resp)

	resp, err = cli.GetCampaignById(apiRequestCtx, campaignTestCase.Id)
	require.NoError(t, err)
	t.Log(resp)

	resp, err = cli.ReplaceCampaign(apiRequestCtx, campaignTestCase, "")
	require.NoError(t, err)
	t.Log(resp)

	ok, err := cli.DeleteCampaign(apiRequestCtx, campaignTestCase.Id)
	require.NoError(t, err)
	t.Log("token context deleted? ", ok)
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
