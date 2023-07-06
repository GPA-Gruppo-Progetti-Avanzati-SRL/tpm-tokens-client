package tokensclient_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/hartracing"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/hartracing/filetracer"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-client/restclient"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-tokens-client/tokensclient"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-tokens-client/tokensclient/model/bearer"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-tokens-client/tokensclient/model/token"
	"net/http"

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

var cliConfig = tokensclient.Config{
	Config: restclient.Config{
		RestTimeout:      0,
		SkipVerify:       false,
		Headers:          nil,
		TraceGroupName:   "tpm-tokens-client-api",
		TraceRequestName: "",
		RetryCount:       0,
		RetryWaitTime:    0,
		RetryMaxWaitTime: 0,
		RetryOnHttpError: nil,
	},
	Host: tokensclient.HostInfo{
		Scheme:   "http",
		HostName: "localhost",
		Port:     8081,
	},
}

func TestMain(m *testing.M) {
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestTokenContextClient(t *testing.T) {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	c, err := InitTracing(t)
	require.NoError(t, err)
	if c != nil {
		defer c.Close()
	}

	cli, err := tokensclient.NewTokensApiClient(&cliConfig)
	require.NoError(t, err)
	defer cli.Close()

	executeTestTokenContextClient(t, cli, &tokenContextTestCase001)
}

func executeTestTokenContextClient(t *testing.T, cli *tokensclient.Client, tokenContextTestCase *token.TokenContext) {
	apiRequestCtx := tokensclient.NewApiRequestContext(tokensclient.ApiRequestWithApiKey("ApiKeyTpmTokens"))

	tokenContextTestCase.Id = "TESTCTX"
	resp, err := cli.NewTokenContext(apiRequestCtx, tokenContextTestCase, "")
	require.NoError(t, err)
	t.Log(resp)

	resp, err = cli.GetTokenContextById(apiRequestCtx, tokenContextTestCase.Id)
	require.NoError(t, err)
	t.Log(resp)

	resp, err = cli.ReplaceTokenContext(apiRequestCtx, tokenContextTestCase, "")
	require.NoError(t, err)
	t.Log(resp)

	ok, err := cli.DeleteTokenContext(apiRequestCtx, tokenContextTestCase.Id)
	require.NoError(t, err)
	t.Log("token context deleted? ", ok)
}

func TestTokenClient(t *testing.T) {
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

	harTracingSpan := hartracing.GlobalTracer().StartSpan()
	defer harTracingSpan.Finish()

	cli, err := tokensclient.NewTokensApiClient(&cliConfig, restclient.WithHarSpan(harTracingSpan))
	require.NoError(t, err)
	defer cli.Close()

	executeTestTokenClient(t, cli, &tokenContextTestCase001)
}

func executeTestTokenClient(t *testing.T, cli *tokensclient.Client, tokenContextTestCase *token.TokenContext) {
	tokenContextTestCase.Id = "BPMGM1"
	apiRequestNew := tokensclient.TokenApiRequest{
		CustomData: map[string]interface{}{
			"ssn":     "MPRMLS62S21G337J",
			"channel": "UP",
			"product": "CC",
		},
	}
	apiRequestCtx := tokensclient.NewApiRequestContext(tokensclient.ApiRequestWithApiKey("ApiKeyLocal"), tokensclient.ApiRequestWithAutoLraId())
	resp, err := cli.NewToken(apiRequestCtx, tokenContextTestCase.Id, &apiRequestNew, "")
	require.NoError(t, err)
	t.Logf("new-token [%s] - NoEvents: %d - json: %s", resp.Id, len(resp.Events), string(resp.MustToJSON()))

	tokenId := resp.Id
	defer func() {
		apiRequestCtx = tokensclient.NewApiRequestContext(tokensclient.ApiRequestWithApiKey("ApiKeyLocal"))
		apiResp, err := cli.DeleteToken(apiRequestCtx, tokenContextTestCase.Id, tokenId)
		require.NoError(t, err)
		t.Log("delete-token - " + string(apiResp.ToJSON()))
	}()

	apiRequestCtx = tokensclient.NewApiRequestContext(tokensclient.ApiRequestWithApiKey("ApiKeyLocal"), tokensclient.ApiRequestWithLraId(apiRequestCtx.LRAId))
	resp, err = cli.CommitToken(apiRequestCtx, tokenContextTestCase.Id, tokenId)
	require.NoError(t, err)
	t.Logf("token-commit [%s] - NoEvents: %d - json: %s", resp.Id, len(resp.Events), string(resp.MustToJSON()))

	apiRequestNextGenerated2Valid1 := tokensclient.TokenApiRequest{
		CustomData: map[string]interface{}{
			"ssn":     "MPRMLS62S21G337J",
			"channel": "UP",
			"product": "CC",
		},
	}
	apiRequestCtx = tokensclient.NewApiRequestContext(tokensclient.ApiRequestWithApiKey("ApiKeyLocal"), tokensclient.ApiRequestWithAutoLraId())
	resp, err = cli.TokenNext(apiRequestCtx, tokenContextTestCase.Id, tokenId, &apiRequestNextGenerated2Valid1, "")
	require.NoError(t, err)
	t.Logf("token-next [%s] - NoEvents: %d - json: %s", resp.Id, len(resp.Events), string(resp.MustToJSON()))

	apiRequestCtx = tokensclient.NewApiRequestContext(tokensclient.ApiRequestWithApiKey("ApiKeyLocal"), tokensclient.ApiRequestWithLraId(apiRequestCtx.LRAId))
	resp, err = cli.CommitToken(apiRequestCtx, tokenContextTestCase.Id, tokenId)
	require.NoError(t, err)
	t.Logf("token-commit [%s] - NoEvents: %d - json: %s", resp.Id, len(resp.Events), string(resp.MustToJSON()))

	apiRequestNextValid12Active := tokensclient.TokenApiRequest{
		CustomData: map[string]interface{}{
			"result": "OK",
		},
	}
	apiRequestCtx = tokensclient.NewApiRequestContext(tokensclient.ApiRequestWithApiKey("ApiKeyLocal"), tokensclient.ApiRequestWithAutoLraId())
	resp, err = cli.TokenNext(apiRequestCtx, tokenContextTestCase.Id, tokenId, &apiRequestNextValid12Active, "")
	require.NoError(t, err)
	t.Logf("token-next [%s] - NoEvents: %d - json: %s", resp.Id, len(resp.Events), string(resp.MustToJSON()))

	apiRequestCtx = tokensclient.NewApiRequestContext(tokensclient.ApiRequestWithApiKey("ApiKeyLocal"), tokensclient.ApiRequestWithLraId(apiRequestCtx.LRAId))
	resp, err = cli.RollbackToken(apiRequestCtx, tokenContextTestCase.Id, tokenId)
	require.NoError(t, err)
	t.Logf("token-rollback [%s] - NoEvents: %d - json: %s", resp.Id, len(resp.Events), string(resp.MustToJSON()))
}

/*
 *
 */

var bearerTestCase001 = bearer.Bearer{
	Pkey:           "MINNIE",
	TokenContextId: "BPMGM1",
	Properties:     nil,
}

func TestBearerClient(t *testing.T) {

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

	harTracingSpan := hartracing.GlobalTracer().StartSpan()
	defer harTracingSpan.Finish()

	cli, err := tokensclient.NewTokensApiClient(&cliConfig, restclient.WithHarTracingEnabled(true), restclient.WithHarSpan(harTracingSpan))
	require.NoError(t, err)
	defer cli.Close()

	executeTestBearerClient(t, cli, &bearerTestCase001)
}

func executeTestBearerClient(t *testing.T, cli *tokensclient.Client, tokenContextTestCase *bearer.Bearer) {
	apiRequestCtx := tokensclient.NewApiRequestContext(
		tokensclient.ApiRequestWithApiKey("ApiKeyTpmTokens"),
		tokensclient.ApiRequestWithHeader("X-Header-Name", "X-Header-Value"))

	bearerRequest := tokensclient.BearerApiRequest{
		Origin: "test",
		TTL:    -1,
	}

	ber, err := cli.AddBearer2Context(apiRequestCtx, bearerTestCase001.Pkey, bearerTestCase001.TokenContextId, &bearerRequest, "")
	require.NoError(t, err)
	t.Log(ber)

	// Remove an header from here
	apiRequestCtx = tokensclient.NewApiRequestContext(
		tokensclient.ApiRequestWithApiKey("ApiKeyTpmTokens"),
	)

	bearerRequest.Properties = map[string]interface{}{"first-name": "PAOLINO"}
	ber, err = cli.UpdateBearerInContext(apiRequestCtx, bearerTestCase001.Pkey, bearerTestCase001.TokenContextId, &bearerRequest, "")
	require.NoError(t, err)
	t.Log(ber)

	resp, err := cli.AddToken2BearerInContext(apiRequestCtx, bearerTestCase001.Pkey, bearerTestCase001.TokenContextId, "TOKEN-ID", "secondary")
	require.NoError(t, err)
	t.Log(resp)

	resp, err = cli.GetBearerInContext(apiRequestCtx, bearerTestCase001.Pkey, bearerTestCase001.TokenContextId)
	require.NoError(t, err)
	t.Log(resp)

	resp, err = cli.RemoveTokenFromBearerInContext(apiRequestCtx, bearerTestCase001.Pkey, bearerTestCase001.TokenContextId, "TOKEN-ID")
	require.NoError(t, err)
	t.Log(resp)

	ber, err = cli.RemoveBearerFromContext(apiRequestCtx, bearerTestCase001.Pkey, bearerTestCase001.TokenContextId)
	require.NoError(t, err)
	t.Log(resp)

	resp, err = cli.GetBearerInContext(apiRequestCtx, bearerTestCase001.Pkey, bearerTestCase001.TokenContextId)
	if err != nil {
		if terr, ok := err.(*tokensclient.ApiResponse); ok {
			if terr.StatusCode != http.StatusNotFound {
				t.Fatal(err)
			}
		}
	}

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
