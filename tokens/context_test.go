package tokens_test

import (
	_ "embed"
	"encoding/json"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-tokens-model/tokens"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	executeTestContext(t, &tokenContextTestCase001)
}

func executeTestContext(t *testing.T, tokenContextTestCase *tokens.TokenContext) {

	tokenContextTestCase.PostProcess()

	ctxJson, err := json.Marshal(tokenContextTestCase)
	require.NoError(t, err)

	t.Log(string(ctxJson))

	tok, err := tokenContextTestCase.NewToken("001", map[string]interface{}{"ssn": "an actor id", "channel": "a channel", "product": "a product"})
	require.NoError(t, err)
	logToken(t, tok)

	time.Sleep(1 * time.Second)

	tok, err = tokenContextTestCase.TokenNext("002", tok, map[string]interface{}{"ssn": "an actor id", "channel": "a channel", "product": "a product"})
	require.NoError(t, err)
	logToken(t, tok)

	time.Sleep(1 * time.Second)

	tok, err = tokenContextTestCase.TokenNext("003", tok, map[string]interface{}{"result": "OK"})
	require.NoError(t, err)
	logToken(t, tok)

	time.Sleep(1 * time.Second)

	tok, err = tokenContextTestCase.TokenNext("003", tok, map[string]interface{}{"result": "OK"})
	require.Equal(t, tokens.RequestIdAlreadyProcessed, err, "Should be idempotent....")
	logToken(t, tok)
}

func logToken(t *testing.T, token *tokens.Token) {
	b, err := json.Marshal(token)
	require.NoError(t, err)
	t.Log(string(b))
}
