package tokens

import (
	"errors"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/expression"
	"github.com/rs/zerolog/log"
	"time"
)

var RequestIdAlreadyProcessed = &TokError{
	Code: TokenDupRequestError,
	Text: MapErrorCode(TokenDupRequestError).Text,
}

type TokenActionOptions struct {
	EventType     EventType
	TxTypePending bool
	LRAId         string
}

type NewTokenOptions struct {
	TokenActionOptions
	TokenId         string
	TokenIdProvider TokenIdProvider
}

type NewTokenOption func(opts *NewTokenOptions)

func NewTokenWithLRAId(lraId string) NewTokenOption {
	return func(opts *NewTokenOptions) {
		if lraId != "" {
			opts.TxTypePending = true
			opts.LRAId = lraId
		}
	}
}

func NewTokenWithId(id string) NewTokenOption {
	return func(opts *NewTokenOptions) {
		opts.TokenId = id
	}
}

func NewTokenWithEventType(evt EventType) NewTokenOption {
	return func(opts *NewTokenOptions) {
		opts.EventType = evt
	}
}

func NewTokenWithTokenIdProvider(provider TokenIdProvider) NewTokenOption {
	return func(opts *NewTokenOptions) {
		if provider != nil {
			opts.TokenIdProvider = provider
		}
	}
}

type NextEventOptions struct {
	TokenActionOptions
}

type NextEventOption func(opts *NextEventOptions)

func NextWithWithLRAId(lraId string) NextEventOption {
	return func(opts *NextEventOptions) {
		if lraId != "" {
			opts.TxTypePending = true
			opts.LRAId = lraId
		}
	}
}

func NextWithEventType(evt EventType) NextEventOption {
	return func(opts *NextEventOptions) {
		opts.EventType = evt
	}
}

func (ctx *TokenContext) NewToken(reqId string, params map[string]interface{}, opt ...NewTokenOption) (*Token, error) {

	const semLogContext = "token-context::new-token"

	if !ctx.IsActive() {
		return nil, NewError(TokenContextNotActiveError, "")
	}

	opts := NewTokenOptions{TokenIdProvider: &DefaultTokenIdProvider{}}
	for _, o := range opt {
		o(&opts)
	}

	// Create  a new context for expression evaluation
	eCtx, err := expression.NewContext(expression.WithMapInput(params))
	if err != nil {
		return nil, NewError(TokenErrorSystemConfiguration, err.Error())
	}

	// Determine the transition to take from state and context (rules and process vars)
	t, err := ctx.StateMachine.selectTransitionFromState(StartEndState, eCtx)
	if err != nil {
		return nil, err
	}

	log.Trace().Str("from", StartEndState).Str("to", t.To).Msg(semLogContext + " found transition")

	// Find the definition of the to state
	var sd StateDefinition
	sd, err = ctx.StateMachine.FindStateDefinition(t.To)
	if err != nil {
		return nil, err
	}

	// Compute the new process vars
	expTs := computeExpiryTimestamp(ctx.Timeline, t.TTL, "")
	pv, err := t.EvalProcessVars(eCtx, expTs)
	if err != nil {
		return nil, err
	}

	// Event type can be overridden by options
	eventType := opts.EventType
	if eventType == "" {
		eventType = EventTypeCreate
	}

	evt := Event{
		RequestId:        reqId,
		EventType:        eventType,
		EventDescription: t.Description,
		State:            State{Code: sd.Code, Description: sd.Description, Pending: opts.TxTypePending, LRAId: opts.LRAId},
		Ts:               time.Now().Format(time.RFC3339),
		Vars:             pv,
		ExpiryTs:         expTs,
	}

	tokId := opts.TokenId
	if tokId == "" {
		act, err := EvaluateFirstActionDefinition(t.Actions, eCtx, ActionTypeNewId)
		if err != nil {
			return nil, err
		}
		tokId, err = opts.TokenIdProvider.NewId(ctx.Id, false, act)
		if err != nil {
			return nil, NewError(TokenErrorNewTokenId, err.Error())
		}
	}

	evt.Actions, err = EvaluateActionDefinitions(t.Actions, eCtx, ActionTypeOut, false)
	if err != nil {
		return nil, NewError(TokenErrorExpressionEvaluation, err.Error())
	}

	tok := Token{
		Pkey:     tokId,
		Id:       tokId,
		Events:   []Event{evt},
		Metadata: nil,
	}

	return &tok, nil
}

func (ctx *TokenContext) TokenNext(reqId string, tok *Token, params map[string]interface{}, opt ...NextEventOption) (*Token, error) {

	const semLogContext = "token-context::token-next"

	if !ctx.IsActive() {
		return nil, NewError(TokenContextNotActiveError, "")
	}

	if tok == nil || len(tok.Events) == 0 {
		return nil, errors.New("invalid token provided")
	}

	tokenCurrentState := tok.FindCurrentState()
	if tokenCurrentState == "" {
		return nil, NewError(TokenErrorInvalidState, "")
	}

	tokenCurrentStateDefinition, err := ctx.StateMachine.FindStateDefinition(tokenCurrentState)
	if err != nil {
		return nil, NewError(TokenErrorContextDefinition, err.Error())
	}

	// Check if request has already been processed
	ndx := tok.FindEventIndexByRequestId(reqId)
	if ndx >= 0 {
		return tok, RequestIdAlreadyProcessed
	}

	switch tokenCurrentStateDefinition.StateType {
	case StateExpired:
		return nil, NewError(TokenExpiredError, tokenCurrentStateDefinition.Help)
	case StateFinal:
		return nil, NewError(TokenFinalStateAlreadyReachedError, tokenCurrentStateDefinition.Help)
	}

	if tok.IsExpired(ctx.Timeline.ExpirationMode) {
		expSd, _ := ctx.StateMachine.FindStateDefinitionByType(StateExpired)
		tok.Events = append(tok.Events, newExpiredEvent(reqId, expSd))
		return tok, NewError(TokenExpiredError, "")
	}

	// Check if token is in pending state
	lastEvt := tok.Events[len(tok.Events)-1]
	if lastEvt.State.Pending {
		return nil, NewError(TokenErrorTransactionInvalidState, fmt.Sprintf("token %s is in pending state and is locked", tok.Id))
	}

	// Create  a new context for expression evaluation
	eCtx, err := expression.NewContext(expression.WithVars(lastEvt.Vars), expression.WithMapInput(params))
	if err != nil {
		return nil, err
	}

	// Determine the transition to take from state and context (rules and process vars)
	t, err := ctx.StateMachine.selectTransitionFromState(tokenCurrentState, eCtx)
	if err != nil {
		return nil, err
	}

	log.Trace().Str("from", tokenCurrentState).Str("to", t.To).Msg(semLogContext + " found transition")

	// Find the definition of the to state
	// variable sd has been reused
	tokenNextStateDefinition, err := ctx.StateMachine.FindStateDefinition(t.To)
	if err != nil {
		return nil, err
	}

	// Determine if the state has already been reached previously. Used to compute the expiry ts.
	ndx = tok.FindLastEventIndex()
	var currentExpiryTs string
	if ndx >= 0 {
		currentExpiryTs = tok.Events[ndx].ExpiryTs
	}

	// Compute the new process vars. final state cannot expire...
	expTs := ""
	if tokenNextStateDefinition.StateType != StateFinal {
		expTs = computeExpiryTimestamp(ctx.Timeline, t.TTL, currentExpiryTs)
	}
	pv, err := t.EvalProcessVars(eCtx, expTs)
	if err != nil {
		return nil, err
	}

	// Merge new variables with previous vars
	if pv == nil {
		pv = make(ProcessVars)
	}
	pv.Merge(lastEvt.Vars)

	opts := NextEventOptions{}
	for _, o := range opt {
		o(&opts)
	}

	// Event type can be overridden by options
	eventType := opts.EventType
	if eventType == "" {
		eventType = EventTypeNext
	}

	evt := Event{
		RequestId:        reqId,
		EventType:        eventType,
		EventDescription: t.Description,
		State:            State{Code: tokenNextStateDefinition.Code, Description: tokenNextStateDefinition.Description, Pending: opts.TxTypePending, LRAId: opts.LRAId},
		Ts:               time.Now().Format(time.RFC3339),
		Vars:             pv,
		ExpiryTs:         expTs,
	}

	evt.Actions, err = EvaluateActionDefinitions(t.Actions, eCtx, ActionTypeOut, false)
	if err != nil {
		return nil, NewError(TokenErrorExpressionEvaluation, err.Error())
	}

	tok.Events = append(tok.Events, evt)

	return tok, nil
}

func computeExpiryTimestamp(ctxTimeline Timeline, ttlDef TTLDefinition, currentExpiryTs string) string {

	const semLogContext = "compute-expiry-timestamp"

	var expTs time.Time
	// TTL has been set....
	if ttlDef.Value != "" {
		ttl := util.ParseDuration(ttlDef.Value, time.Hour*24)
		expTs = time.Now().Add(ttl)
		if ctxTimeline.ExpirationMode == ExpirationModeTimestamp {
			return expTs.Format(time.RFC3339)
		}

		return expTs.Format("20060102")
	}

	if currentExpiryTs == "" {
		if ctxTimeline.ExpirationMode == "" || ctxTimeline.ExpirationMode == ExpirationModeDate {
			return ctxTimeline.EndDate
		}

		tm, err := time.Parse("20060102", ctxTimeline.EndDate)
		if err != nil {
			log.Error().Str("start-date", ctxTimeline.StartDate).Str("end-date", ctxTimeline.EndDate).Err(err).Msg(semLogContext + " context timeline invalid format")
		}

		return tm.Add(time.Hour * 23).Add(time.Minute * 59).Add(time.Second * 59).Format(time.RFC3339)
	}

	return currentExpiryTs
}

func (ctx *TokenContext) TokenCommit(reqId string, tok *Token, opt ...NextEventOption) (*Token, bool, error) {

	const semLogContext = "token-context::commit"

	if tok == nil || len(tok.Events) == 0 {
		return nil, false, NewError(TokenErrorTransactionInvalidState, "no events found in the chain")
	}

	opts := NextEventOptions{}
	for _, o := range opt {
		o(&opts)
	}

	// Check if request has already been processed
	ndx, err := tok.FindEventIndexInPendingState4Commit(opts.LRAId)
	if err != nil {
		return nil, false, err
	}

	done := false
	if ndx >= 0 && tok.Events[ndx].State.Pending {
		log.Info().Int("evt-ndx", ndx).Msg(semLogContext + " doing commit of event")
		tok.Events[ndx].State.Pending = false
		tok.Events = append(tok.Events, dupEventOf(reqId, &tok.Events[ndx], EventTypeCommit, true))
		done = true
	}

	return tok, done, nil
}

func (ctx *TokenContext) TokenRollback(reqId string, tok *Token, opt ...NextEventOption) (*Token, bool, error) {

	const semLogContext = "token-context::rollback"

	if tok == nil || len(tok.Events) == 0 {
		return nil, false, NewError(TokenErrorTransactionInvalidState, "no events found in the chain")
	}

	opts := NextEventOptions{}
	for _, o := range opt {
		o(&opts)
	}

	// Check if request has already been processed
	ndx, err := tok.FindEventIndexInPendingState4Rollback(opts.LRAId)
	if err != nil {
		return nil, false, err
	}

	done := false
	if ndx >= 0 {
		log.Info().Int("evt-ndx", ndx).Msg(semLogContext + " doing rollback of event")
		var rollbackEvent Event
		tok.Events[ndx].State.Pending = false
		if ndx > 0 {
			rollbackEvent = dupEventOf(reqId, &tok.Events[ndx-1], EventTypeRollback, false)
		} else {
			// I'm rollbacking the last event.... it's an empty token that should be handled....
			rollbackEvent = dupEventOf(reqId, &tok.Events[ndx], EventTypeRollback, false)
			rollbackEvent.State.Code = StartEndState
		}
		tok.Events[ndx] = rollbackEvent

		done = true
	} else {
		log.Info().Msg(semLogContext + " no event to rollback found")
	}

	return tok, done, nil
}

func dupEventOf(reqId string, evt *Event, eventType EventType, dupActions bool) Event {
	newEvent := Event{
		RequestId:        reqId,
		EventType:        eventType,
		EventDescription: evt.EventDescription,
		State:            evt.State,
		Ts:               time.Now().Format(time.RFC3339),
		ExpiryTs:         evt.ExpiryTs,
		Vars:             nil,
		Actions:          nil,
	}

	if len(evt.Vars) > 0 {
		newEvent.Vars = make(ProcessVars)
		for n, v := range evt.Vars {
			newEvent.Vars[n] = v
		}
	}

	if dupActions && len(evt.Actions) > 0 {
		for _, a := range evt.Actions {
			newEvent.Actions = append(newEvent.Actions, a)
		}
	}
	return newEvent
}

func newExpiredEvent(reqId string, expiredState StateDefinition) Event {
	newEvent := Event{
		RequestId: reqId,
		EventType: EventTypeExpiration,
		State: State{
			Code:        expiredState.Code,
			Description: expiredState.Description,
		},
		Ts: time.Now().Format(time.RFC3339),
	}

	return newEvent
}
