package tokens

import (
	"github.com/rs/zerolog/log"
)

func (pvs ProcessVars) Merge(pvs2 ProcessVars) {
	for n, v := range pvs2 {
		if _, ok := pvs[n]; !ok {
			pvs[n] = v
		}
	}
}

func (tok *Token) FindEventIndexInPendingState4Commit(lraId string) (int, error) {

	if len(tok.Events) == 0 {
		if lraId != "" {
			return -1, NewError(TokenErrorTransactionInvalidState, "could not find an LRAId committed or not.")
		}
		return -1, nil
	}

	lraIdNdx := -1
	pendingNdx := -1
	for i := len(tok.Events) - 1; i >= 0; i-- {
		e := tok.Events[i]
		if e.State.Pending {
			if pendingNdx >= 0 {
				log.Error().Msg("multiple pending states found...")
			} else {
				pendingNdx = i
			}
		}
		if lraId != "" && lraId == e.State.LRAId {
			if lraIdNdx >= 0 {
				log.Error().Msg("multiple lraIds found...")
			} else {
				lraIdNdx = i
			}
		}
	}

	log.Trace().Int("pending-evt-ndx", pendingNdx).Int("lraid-evt-ndx", lraIdNdx).Str("lra-id", lraId).Msg("last pending event")
	if lraId == "" {
		if pendingNdx >= 0 && pendingNdx != (len(tok.Events)-1) {
			return -1, NewError(TokenErrorTransactionInvalidState, "the committable event is not the last event")
		}
		return pendingNdx, nil
	}

	if lraIdNdx < 0 {
		return lraIdNdx, NewError(TokenErrorTransactionInvalidState, "could not find an LRAId committed or not.")
	}

	if lraIdNdx >= 0 && lraIdNdx != (len(tok.Events)-1) {
		return -1, NewError(TokenErrorTransactionInvalidState, "the committable event is not the last event")
	}

	return lraIdNdx, nil
}

func (tok *Token) FindEventIndexInPendingState4Rollback(lraId string) (int, error) {

	if len(tok.Events) == 0 {
		return -1, nil
	}

	evt := tok.Events[len(tok.Events)-1]
	if evt.State.Pending {
		if lraId != "" {
			if evt.State.LRAId != lraId {
				return -1, NewError(TokenErrorTransactionInvalidState, "inconsistent operation... lraid doesn't match the pending lraid")
			}
		}

		return len(tok.Events) - 1, nil
	}

	return -1, nil
}
