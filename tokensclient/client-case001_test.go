package tokensclient_test

import "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-tokens-client/tokensclient"

// This file is being separated from the actual tests because in the current state of affairs it has to be aligned with the tokensazstore package one and the client one.
// easier to copy a file and change package instead of copying and pasting

var tokenContextTestCase001 = tokensclient.TokenContext{
	Id:        "BPMIFI",
	Pkey:      tokensclient.ContextPartitionKey,
	Platform:  "BP",
	Version:   tokensclient.TokenContextBaseVersion,
	Suspended: false,
	Timeline: tokensclient.Timeline{
		StartDate:      "20230101",
		EndDate:        "20230330",
		ExpirationMode: tokensclient.ExpirationModeDate,
	},
	StateMachine: tokensclient.StateMachine{
		States: []tokensclient.StateDefinition{
			{
				Code:        "[*]",
				StateType:   tokensclient.StartEndState,
				Description: "description of start/end meta-state",
				Help:        "help of [*] state",
				Actions: []tokensclient.ActionDefinition{
					{
						ActionId:   "domain-specific-in-action-id",
						ActionType: tokensclient.ActionTypeIn,
						Properties: map[string]interface{}{
							"cf":      "{$.ssn}",
							"product": "{$.product}",
							"channel": "{$.channel}",
						},
					},
				},
				OutTransitions: []tokensclient.Transition{
					{
						To:          "generated",
						Description: "il codice e' stato creato correttamente",
						Properties: []tokensclient.Property{
							{
								Name:           "ssn",
								ValidationRule: "required",
							},
							{
								Name:           "channel",
								ValidationRule: "required",
							},
							{
								Name:           "product",
								ValidationRule: "required",
							},
						},
						ProcessVarDefinitions: []tokensclient.ProcessVarDefinition{
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
						TTL: tokensclient.TTLDefinition{
							Value: "1d",
						},
						Actions: []tokensclient.ActionDefinition{
							{
								ActionId:   tokensclient.ActionTypeNewId,
								ActionType: tokensclient.ActionTypeNewId,
								Properties: map[string]interface{}{
									"cf": "{$.ssn}",
								},
							},
							{
								ActionId:   "domain-specific-out-action-id",
								ActionType: tokensclient.ActionTypeOut,
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
				StateType:   tokensclient.StateStd,
				Description: "description of generated state",
				Help:        "help of generated state",
				BusinessView: tokensclient.BusinessViewState{
					Code:        "whatever",
					Description: "whatever description",
				},
				OutTransitions: []tokensclient.Transition{
					{
						To:          "valid-1",
						Description: "il codice e' stato usato ed in corso di validazione",
						Properties: []tokensclient.Property{
							{
								Name:           "ssn",
								ValidationRule: "required",
							},
							{
								Name:           "channel",
								ValidationRule: "required",
							},
							{
								Name:           "product",
								ValidationRule: "required",
							},
						},
						Rules: []tokensclient.Rule{
							{
								Expression: "\"{$.ssn}\" == \"{v:ssn1}\" && \"{$.channel}\" == \"{v:channel}\" && \"{$.product}\" == \"{v:product}\"",
							},
						},
						ProcessVarDefinitions: nil,
						TTL: tokensclient.TTLDefinition{
							Value: "1d",
						},
					},
				},
			},
			{
				Code:        "valid-1",
				StateType:   tokensclient.StateStd,
				Description: "description of valid-1 state",
				Help:        "help of valid-1 state",
				OutTransitions: []tokensclient.Transition{
					{
						To: "active",
						Properties: []tokensclient.Property{
							{
								Name:           "result",
								ValidationRule: "required",
							},
						},
						Rules: []tokensclient.Rule{
							{
								Expression: "\"{$.result}\" == \"OK\"",
							},
						},
					},
					{
						To: "generated",
						Properties: []tokensclient.Property{
							{
								Name:           "result",
								ValidationRule: "required",
							},
						},
						Rules: []tokensclient.Rule{
							{
								Expression: "\"{$.result}\" == \"KO\"",
							},
						},
					},
				},
			},
			{
				Code:        "active",
				StateType:   tokensclient.StateStd,
				Description: "description of active state",
				Help:        "help of active state",
				OutTransitions: []tokensclient.Transition{
					{
						To: "valid-2",
						Properties: []tokensclient.Property{
							{
								Name:           "ssn",
								ValidationRule: "required",
							},
							{
								Name:           "channel",
								ValidationRule: "required",
							},
							{
								Name:           "product",
								ValidationRule: "required",
							},
						},
						Rules: []tokensclient.Rule{
							{
								Expression: "\"{$.ssn}\" != \"{v:ssn}\" && \"{$.channel}\" == \"UP\"",
							},
						},
					},
				},
			},
			{
				Code:        "valid-2",
				StateType:   tokensclient.StateStd,
				Description: "description of valid-2 state",
				Help:        "help of valid-2 state",
				OutTransitions: []tokensclient.Transition{
					{
						To: "consumed",
						Properties: []tokensclient.Property{
							{
								Name:           "result",
								ValidationRule: "required",
							},
						},
						Rules: []tokensclient.Rule{
							{
								Expression: "\"{$.result}\" == \"OK\"",
							},
						},
					},
				},
			},
			{
				Code:         "consumed",
				StateType:    tokensclient.StateFinal,
				Description:  "description of consumed state",
				Help:         "help of consumed state",
				BusinessView: tokensclient.BusinessViewState{},
			},
			{
				Code:        "expired",
				StateType:   tokensclient.StateExpired,
				Description: "description of expired state",
				Help:        "help of expired state",
			},
		},
	},
	TokenIdProviderType: &tokensclient.TokenIdProviderType{ProviderType: tokensclient.TokenIdProviderTypeExternal},
}
