package tokens_test

import "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-tokens-model/tokens"

// This file is being separated from the actual tests because in the current state of affairs it has to be aligned with the tokensazstore package one and the client one.
// easier to copy a file and change package instead of copying and pasting

var tokenContextTestCase001 = tokens.TokenContext{
	Id:        "BPMIFI",
	Pkey:      tokens.ContextPartitionKey,
	Platform:  "BP",
	Version:   tokens.TokenContextBaseVersion,
	Suspended: false,
	Timeline: tokens.Timeline{
		StartDate:      "20230101",
		EndDate:        "20230330",
		ExpirationMode: tokens.ExpirationModeDate,
	},
	StateMachine: tokens.StateMachine{
		States: []tokens.StateDefinition{
			{
				Code:        "[*]",
				StateType:   tokens.StartEndState,
				Description: "description of start/end meta-state",
				Help:        "help of [*] state",
				Actions: []tokens.ActionDefinition{
					{
						ActionId:   "domain-specific-in-action-id",
						ActionType: tokens.ActionTypeIn,
						Properties: map[string]interface{}{
							"cf":      "{$.ssn}",
							"product": "{$.product}",
							"channel": "{$.channel}",
						},
					},
				},
				OutTransitions: []tokens.Transition{
					{
						To:          "generated",
						Description: "il codice e' stato creato correttamente",
						Rules: []tokens.Rule{
							{
								Expression: "\"{$.ssn}\" != \"\"",
							},
							{
								Expression: "\"{$.channel}\" != \"\"",
							},
							{
								Expression: "\"{$.product}\" != \"\"",
							},
						},
						ProcessVarDefinitions: []tokens.ProcessVarDefinition{
							{
								Name:        "ssn1",
								Description: "customer id: social security number",
								Value:       "{$.ssn}",
							},
							{
								Name:        "channel",
								Description: "a generic property used for routing",
								Value:       "{$.channel}",
							},
							{
								Name:        "product",
								Description: "code of a product",
								Value:       "{$.product}",
							},
						},
						TTL: tokens.TTLDefinition{
							Value: "1d",
						},
						Actions: []tokens.ActionDefinition{
							{
								ActionId:   tokens.ActionTypeNewId,
								ActionType: tokens.ActionTypeNewId,
								Properties: map[string]interface{}{
									"cf": "{$.ssn}",
								},
							},
							{
								ActionId:   "domain-specific-out-action-id",
								ActionType: tokens.ActionTypeOut,
								Properties: map[string]interface{}{
									"cf":      "{$.ssn}",
									"product": "{v:product}",
								},
							},
						},
					},
				},
			},
			{
				Code:        "generated",
				StateType:   tokens.StateStd,
				Description: "description of generated state",
				Help:        "help of generated state",
				BusinessView: tokens.BusinessViewState{
					Code:        "whatever",
					Description: "whatever description",
				},
				OutTransitions: []tokens.Transition{
					{
						To:          "valid-1",
						Description: "il codice e' stato usato ed in corso di validazione",
						Rules: []tokens.Rule{
							{
								Expression: "\"{$.ssn}\" == \"{v:ssn1}\" && \"{$.channel}\" == \"{v:channel}\" && \"{$.product}\" == \"{v:product}\"",
							},
						},
						ProcessVarDefinitions: nil,
						TTL: tokens.TTLDefinition{
							Value: "1d",
						},
					},
				},
			},
			{
				Code:        "valid-1",
				StateType:   tokens.StateStd,
				Description: "description of valid-1 state",
				Help:        "help of valid-1 state",
				OutTransitions: []tokens.Transition{
					{
						To: "active",
						Rules: []tokens.Rule{
							{
								Expression: "\"{$.result}\" == \"OK\"",
							},
						},
					},
					{
						To: "generated",
						Rules: []tokens.Rule{
							{
								Expression: "\"{$.result}\" == \"KO\"",
							},
						},
					},
				},
			},
			{
				Code:        "active",
				StateType:   tokens.StateStd,
				Description: "description of active state",
				Help:        "help of active state",
				OutTransitions: []tokens.Transition{
					{
						To: "valid-2",
						Rules: []tokens.Rule{
							{
								Expression: "\"{$.ssn}\" != \"{v:ssn}\" && \"{$.channel}\" == \"UP\"",
							},
						},
					},
				},
			},
			{
				Code:        "valid-2",
				StateType:   tokens.StateStd,
				Description: "description of valid-2 state",
				Help:        "help of valid-2 state",
				OutTransitions: []tokens.Transition{
					{
						To: "consumed",
						Rules: []tokens.Rule{
							{
								Expression: "\"{$.result}\" == \"OK\"",
							},
						},
					},
				},
			},
			{
				Code:         "consumed",
				StateType:    tokens.StateFinal,
				Description:  "description of consumed state",
				Help:         "help of consumed state",
				BusinessView: tokens.BusinessViewState{},
			},
			{
				Code:        "expired",
				StateType:   tokens.StateExpired,
				Description: "description of expired state",
				Help:        "help of expired state",
			},
		},
	},
	TokenIdProviderType: &tokens.TokenIdProviderType{ProviderType: tokens.TokenIdProviderTypeExternal},
}
