package tokens

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/expression"
)

//var TransitionNotFound = errors.New("transition not found")
//var StateDefinitionNotFound = errors.New("state definition not found")

func (sm *StateMachine) FindStateDefinitionByType(st StateType) (StateDefinition, int) {

	foundNdx := -1
	numFounds := 0
	for i, s := range sm.States {
		if s.StateType == st {
			foundNdx = i
			numFounds++
		}
	}

	if foundNdx >= 0 {
		return sm.States[foundNdx], numFounds
	}

	return StateDefinition{}, numFounds
}

func (sm *StateMachine) findOrderedTransitionsFromState(st string) ([]Transition, error) {

	stateDef, err := sm.FindStateDefinition(st)
	if err != nil {
		return nil, err
	}

	return stateDef.OutTransitions, err

	/*
		if len(ts) > 0 {
			sort.SliceStable(ts, func(p, q int) bool {
				if ts[p].Order <= 0 && ts[q].Order <= 0 {
					return false
				}

				if ts[p].Order <= 0 {
					return false
				}

				if ts[q].Order <= 0 {
					return true
				}

				return ts[p].Order < ts[q].Order
			})
		}

		return ts, nil
	*/

}

func (sm *StateMachine) selectTransitionFromState(st string, exprContext *expression.Context) (Transition, error) {

	var err error
	ts, _ := sm.findOrderedTransitionsFromState(st)
	if len(ts) > 0 {
		for _, t := range ts {
			rs := true
			if len(t.Rules) > 0 {
				rs, err = sm.BoolEvalRules(exprContext, t.Rules, expression.AllMustMatch)
			}

			if err != nil {
				return Transition{}, NewError(TokenErrorExpressionEvaluation, err.Error())
			}

			if rs {
				return t, nil
			}
		}
	}

	sd, err := sm.FindStateDefinition(st)
	if err != nil {
		return Transition{}, err
	}

	return Transition{}, NewError(TokenErrorNotTransitionFound, sd.Help)
}

func (sm *StateMachine) BoolEvalRules(exprContext *expression.Context, varExpressions []Rule, mode expression.EvaluationMode) (bool, error) {
	if len(varExpressions) == 0 {
		return false, nil
	}

	exprs := make([]string, len(varExpressions))
	for i, r := range varExpressions {
		exprs[i] = r.Expression
	}

	return exprContext.BoolEvalMany(exprs, mode)
}

func (t *Transition) EvalProcessVars(eCtx *expression.Context, expTs string) (ProcessVars, error) {

	var err error
	if len(t.ProcessVarDefinitions) > 0 || expTs != "" {
		pv := make(ProcessVars)
		for _, pvd := range t.ProcessVarDefinitions {
			pv[pvd.Name], err = eCtx.EvalOne(pvd.Value)
			if err != nil {
				return nil, NewError(TokenErrorExpressionEvaluation, err.Error())
			}

			eCtx.Add(pvd.Name, pv[pvd.Name])
		}

		if expTs != "" {
			pv[ExpiryTsProcessVariable] = expTs
			eCtx.Add(ExpiryTsProcessVariable, expTs)
		}
		return pv, nil
	}

	return nil, nil
}

func EvaluateFirstActionDefinition(actions []ActionDefinition, eCtx *expression.Context, actionType ActionType) (*Action, error) {
	acts, err := EvaluateActionDefinitions(actions, eCtx, actionType, true)
	if err != nil {
		return nil, err
	}

	if len(acts) > 0 {
		a := acts[0]
		return &a, nil
	}

	return nil, nil
}
