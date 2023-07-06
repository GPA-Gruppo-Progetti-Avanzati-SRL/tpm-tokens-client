package tokensclient_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-tokens-client/tokensclient/model/token"
)

// This file is being separated from the actual tests because in the current state of affairs it has to be aligned with the tokensazstore package one and the client one.
// easier to copy a file and change package instead of copying and pasting
const (
	ParamNameCf                   = "cf"
	ParamCfHelpMessage            = "il codice fiscale e' obbligatorio"
	ParamNameCanale               = "canale"
	ParamCanaleHelpMessage        = "il canale e' obbligatorio"
	ParamNameFunnelId             = "funnelId"
	ParamFunnelIdHelpMessage      = "il funnelId e' obbligatorio"
	ParamNameProdotto             = "prodotto"
	ParamProdottoHelpMessage      = "il prodotto e' obbligatorio"
	ParamNameFase                 = "fase"
	ParamFaseHelpMessage          = "la fase e' obbligatoria"
	ParamNameServizio             = "servizio"
	ParamServizioHelpMessage      = "il servizio e' obbligatoria"
	ParamNameNumero               = "numero"
	ParamNumeroHelpMessage        = "il numero conto e' obbligatoria"
	ParamNameNumeroPratica        = "numeroPratica"
	ParamNumeroPraticaHelpMessage = "il numero pratica e' obbligatoria"

	InputParamCf         = "{$." + ParamNameCf + "}"
	InputParamCanale     = "{$." + ParamNameCanale + "}"
	InputParamFunnelId   = "{$." + ParamNameFunnelId + "}"
	InputParamProdotto   = "{$." + ParamNameProdotto + "}"
	InputParamFase       = "{$." + ParamNameFase + "}"
	InputParamServizio   = "{$." + ParamNameServizio + "}"
	InputParamNumero     = "{$." + ParamNameNumero + "}"
	InputParamNumPratica = "{$." + ParamNameNumeroPratica + "}"

	MGMStatusGenerato                     = "generato"
	MGMStatusInAttesaAperturaPrimoConto   = "in-attesa-apertura-primo-conto"
	MGMStatusPrimoContoAperto             = "primo-conto-aperto"
	MGMStatusInAttesaAperturaSecondoConto = "in-attesa-apertura-secondo-conto"
	MGMStatusBruciato                     = "bruciato"
	MGMStatusExpired                      = "expired"

	CodiceNonUtilizzabile = "codice-non-utilizzabile"

	BearerCF1ReferenceVariable = "{v:cf1}"
	BearerCF2ReferenceVariable = "{v:cf2}"
)

var tokenContextTestCase001 = token.TokenContext{
	Id:        "BPMGM1",
	Pkey:      token.ContextPartitionKey,
	Platform:  "BP",
	Version:   token.TokenContextBaseVersion,
	Suspended: false,
	Timeline: token.Timeline{
		StartDate:      "20230101",
		EndDate:        "20230430",
		ExpirationMode: token.ExpirationModeDate,
	},
	StateMachine: token.StateMachine{
		States: []token.StateDefinition{
			{
				Code:        "[*]",
				StateType:   token.StartEndState,
				Description: "description of start/end meta-state",
				/* No out actions are required in here
				Actions: []tokensclient.ActionDefinition{
					{
						ActionId:   "in-action-start-action-id",
						ActionType: tokensclient.ActionTypeIn,
						Properties: map[string]interface{}{
							"cf": InputParamCf,
						},
					},
				},
				*/
				OutTransitions: []token.Transition{
					{
						Name:        "creazione",
						To:          MGMStatusGenerato,
						Description: "il codice e' stato creato per il cf: {v:cf1}",
						Properties: []token.Property{
							{
								Name:           ParamNameCf,
								ValidationRule: "required",
								Help:           token.CodeDescriptionPair{Description: ParamCfHelpMessage},
							},
						},
						Bearers: []token.BearerRef{
							{
								Id:   BearerCF1ReferenceVariable,
								Role: "primary",
							},
						},
						ProcessVarDefinitions: []token.ProcessVarDefinition{
							{
								Name:        "cf1",
								Description: "codice fiscale utilizzatore",
								Value:       InputParamCf,
							},
						},
						TTL: token.TTLDefinition{
							Value: "1d",
						},
						Actions: []token.ActionDefinition{
							{
								ActionId:   token.ActionTypeNewId,
								ActionType: token.ActionTypeNewId,
								Properties: map[string]interface{}{
									ParamNameCf: InputParamCf,
								},
							},
							/* No out actions are required in here
							{
								ActionId:   "domain-specific-out-action-id",
								ActionType: tokensclient.ActionTypeOut,
								Properties: map[string]interface{}{
									"cf": InputParamCf,
								},
							},
							*/
						},
					},
				},
			},
			{
				Code:        MGMStatusGenerato,
				StateType:   token.StateStd,
				Description: "Il codice è in attesa di un primo utilizzo da parte del codice fiscale dell'assegnatario",
				Actions: []token.ActionDefinition{
					{
						ActionId:   "in-action-verifica-15-cf",
						ActionType: token.ActionTypeIn,
						Properties: map[string]interface{}{
							ParamNameCf: InputParamCf,
						},
					},
				},
				/* Not applicable in here
				BusinessView: tokensclient.CodeDescriptionPair{
					Code:        "whatever",
					Description: "whatever description",
				},
				*/
				OutTransitions: []token.Transition{
					{
						Name:        "prenotazione",
						To:          MGMStatusInAttesaAperturaPrimoConto,
						Description: "Il codice e' stato usato da {v:cf1} e il sistema e in attesa di apertura del primo conto",
						Properties: []token.Property{
							{
								Name:           ParamNameCf,
								ValidationRule: "required",
								Help:           token.CodeDescriptionPair{Description: ParamCfHelpMessage},
							},
							{
								Name:           ParamNameCanale,
								ValidationRule: "required",
								Help:           token.CodeDescriptionPair{Description: ParamCanaleHelpMessage},
							},
							{
								Name:           ParamNameProdotto,
								ValidationRule: "required",
								Help:           token.CodeDescriptionPair{Description: ParamProdottoHelpMessage},
							},
							{
								Name:           ParamNameFase,
								ValidationRule: "required",
								Help:           token.CodeDescriptionPair{Description: ParamFaseHelpMessage},
							},
							{
								Name:           ParamNameFunnelId,
								ValidationRule: "required",
								Help:           token.CodeDescriptionPair{Description: ParamFunnelIdHelpMessage},
							},
							{
								Name:           ParamNameNumeroPratica,
								ValidationRule: "required",
								Help:           token.CodeDescriptionPair{Description: ParamNumeroPraticaHelpMessage},
								Scope:          string(token.EventTypeNext),
							},
							{
								Name:           ParamNameServizio,
								ValidationRule: "required",
								Help:           token.CodeDescriptionPair{Description: ParamServizioHelpMessage},
								Scope:          string(token.EventTypeNext),
							},
							{
								Name:           ParamNameNumero,
								ValidationRule: "required",
								Help:           token.CodeDescriptionPair{Description: ParamNumeroHelpMessage},
								Scope:          string(token.EventTypeNext),
							},
						},
						Rules: []token.Rule{
							{
								Expression: "\"{$.cf}\" == \"{v:cf1}\"",
								Help:       token.CodeDescriptionPair{Description: "Il codice inserito non risulta associato al cliente {$.cf} ma al cliente {v:cf1}"},
							},
						},
						Bearers: []token.BearerRef{
							{
								Id:   BearerCF1ReferenceVariable,
								Role: "primary",
							},
						},
						ProcessVarDefinitions: []token.ProcessVarDefinition{
							{
								Name:        "numeroPratica1",
								Description: "numero pratica apertura primo conto",
								Value:       InputParamNumPratica,
							},
							{
								Name:        "servizio1",
								Description: "codice servizio apertura primo conto",
								Value:       InputParamServizio,
							},
							{
								Name:        "numero1",
								Description: "numero primo conto",
								Value:       InputParamNumero,
							},
							{
								Name:        "canale1",
								Description: "canale apertura primo conto",
								Value:       InputParamCanale,
							},
							{
								Name:        "fase1",
								Description: "fase apertura primo conto",
								Value:       InputParamFase,
							},
							{
								Name:        "funnelId1",
								Description: "funnel id apertura primo conto",
								Value:       InputParamFunnelId,
							},
							{
								Name:        "prodotto1",
								Description: "prodotto apertura primo conto",
								Value:       InputParamProdotto,
							},
						},
						TTL: token.TTLDefinition{
							Value: "1d",
						},
					},
				},
			},
			{
				Code:        MGMStatusInAttesaAperturaPrimoConto,
				StateType:   token.StateStd,
				Description: "il sistema e' in attesa del perfezionamento della pratica relativa al primo conto corrente",
				Help: token.CodeDescriptionPair{
					Code:        CodiceNonUtilizzabile,
					Description: "il codice inserito non e' ancora stato attivato",
				},
				OutTransitions: []token.Transition{
					{
						Name:        "attivazione",
						To:          MGMStatusPrimoContoAperto,
						Description: "il perfezionamento della pratica {v:numeroPratica1} di apertura del conto {v:servizio1}-{v:numero1} è terminata con successo",
						Properties: []token.Property{
							{
								Name:           "result",
								ValidationRule: "required",
							},
						},
						Rules: []token.Rule{
							{
								Expression: "\"{$.result}\" == \"OK\"",
							},
						},
						Bearers: []token.BearerRef{
							{
								Id:   BearerCF1ReferenceVariable,
								Role: "primary",
							},
						},
					},
					{
						Name:        "attivazione-fallita",
						To:          MGMStatusGenerato,
						Description: "il perfezionamento della pratica {v:numeroPratica1} di apertura del conto {v:servizio}-{v:numero} non e' andato a buon fine",
						Properties: []token.Property{
							{
								Name:           "result",
								ValidationRule: "required",
							},
						},
						Bearers: []token.BearerRef{
							{
								Id:   BearerCF1ReferenceVariable,
								Role: "primary",
							},
						},
						Rules: []token.Rule{
							{
								Expression: "\"{$.result}\" == \"KO\"",
							},
						},
					},
				},
			},
			{
				Code:        MGMStatusPrimoContoAperto,
				StateType:   token.StateStd,
				Description: "Il primo conto e' stato aperto con successo, il codice puo' essere ora utilizzato da un codice fiscale presentato dall'assegnatario del codice",
				Actions: []token.ActionDefinition{
					{
						ActionId:   "in-action-verifica-15-cf",
						ActionType: token.ActionTypeIn,
						Properties: map[string]interface{}{
							ParamNameCf: InputParamCf,
						},
					},
				},
				OutTransitions: []token.Transition{
					{
						Name: "uso",
						To:   MGMStatusInAttesaAperturaSecondoConto,
						Properties: []token.Property{
							{
								Name:           ParamNameCf,
								ValidationRule: "required",
								Help:           token.CodeDescriptionPair{Description: ParamCfHelpMessage},
							},
							{
								Name:           ParamNameCanale,
								ValidationRule: "required",
								Help:           token.CodeDescriptionPair{Description: ParamCanaleHelpMessage},
							},
							{
								Name:           ParamNameProdotto,
								ValidationRule: "required",
								Help:           token.CodeDescriptionPair{Description: ParamProdottoHelpMessage},
							},
							{
								Name:           ParamNameFase,
								ValidationRule: "required",
								Help:           token.CodeDescriptionPair{Description: ParamFaseHelpMessage},
							},
							{
								Name:           ParamNameFunnelId,
								ValidationRule: "required",
								Help:           token.CodeDescriptionPair{Description: ParamFunnelIdHelpMessage},
							},
							{
								Name:           ParamNameNumeroPratica,
								ValidationRule: "required",
								Help:           token.CodeDescriptionPair{Description: ParamNumeroPraticaHelpMessage},
								Scope:          string(token.EventTypeNext),
							},
							{
								Name:           ParamNameServizio,
								ValidationRule: "required",
								Help:           token.CodeDescriptionPair{Description: ParamServizioHelpMessage},
								Scope:          string(token.EventTypeNext),
							},
							{
								Name:           ParamNameNumero,
								ValidationRule: "required",
								Help:           token.CodeDescriptionPair{Description: ParamNumeroHelpMessage},
								Scope:          string(token.EventTypeNext),
							},
						},
						Bearers: []token.BearerRef{
							{
								Id:   BearerCF1ReferenceVariable,
								Role: "primary",
							},
							{
								Id:   BearerCF2ReferenceVariable,
								Role: "secondary",
							},
						},
						ProcessVarDefinitions: []token.ProcessVarDefinition{
							{
								Name:        "cf2",
								Description: "codice fiscale del presentato",
								Value:       InputParamCf,
							},
							{
								Name:        "numeroPratica2",
								Description: "numero pratica apertura secondo conto",
								Value:       InputParamNumPratica,
							},
							{
								Name:        "servizio2",
								Description: "codice servizio apertura secondo conto",
								Value:       InputParamServizio,
							},
							{
								Name:        "numero2",
								Description: "numero secondo conto",
								Value:       InputParamNumero,
							},
							{
								Name:        "canale2",
								Description: "canale apertura secondo conto",
								Value:       InputParamCanale,
							},
							{
								Name:        "fase2",
								Description: "fase apertura secondo conto",
								Value:       InputParamFase,
							},
							{
								Name:        "funnelId2",
								Description: "funnel id apertura secondo conto",
								Value:       InputParamFunnelId,
							},
							{
								Name:        "prodotto2",
								Description: "prodotto apertura secondo conto",
								Value:       InputParamProdotto,
							},
						},
						Rules: []token.Rule{
							{
								Expression: "\"{$.cf}\" != \"{v:cf1}\"",
								Help: token.CodeDescriptionPair{
									Description: "il codice fiscale dell'utilizzatore {$.cf} deve essere diverso dal codice fiscale dell'assegnatario del codice {v:cf1}",
								},
							},
						},
					},
				},
			},
			{
				Code:        MGMStatusInAttesaAperturaSecondoConto,
				StateType:   token.StateStd,
				Description: "il sistema e' in attesa del perfezionamento della pratica del secondo conto",
				Help: token.CodeDescriptionPair{
					Code:        CodiceNonUtilizzabile,
					Description: "il codice inserito e' stato gia' usato e in attesa di lavorazione",
				},
				OutTransitions: []token.Transition{
					{
						Name: "bruciatura",
						To:   MGMStatusBruciato,
						Properties: []token.Property{
							{
								Name:           "result",
								ValidationRule: "required",
							},
						},
						Bearers: []token.BearerRef{
							{
								Id:   BearerCF1ReferenceVariable,
								Role: "primary",
							},
							{
								Id:   BearerCF2ReferenceVariable,
								Role: "secondary",
							},
						},
						Rules: []token.Rule{
							{
								Expression: "\"{$.result}\" == \"OK\"",
							},
						},
					},
					{
						Name: "bruciatura-fallita",
						To:   MGMStatusPrimoContoAperto,
						Properties: []token.Property{
							{
								Name:           "result",
								ValidationRule: "required",
							},
						},
						Bearers: []token.BearerRef{
							{
								Id:   BearerCF1ReferenceVariable,
								Role: "primary",
							},
						},
						Rules: []token.Rule{
							{
								Expression: "\"{$.result}\" == \"KO\"",
							},
						},
					},
				},
			},
			{
				Code:        MGMStatusBruciato,
				StateType:   token.StateFinal,
				Description: "il codice e' stato utilizzato",
				Help: token.CodeDescriptionPair{
					Code:        CodiceNonUtilizzabile,
					Description: "il codice inserito e' gia' stato utilizzato",
				},
			},
			{
				Code:        MGMStatusExpired,
				StateType:   token.StateExpired,
				Description: "il codice inserito è scaduto e non più utilizzabile",
				Help: token.CodeDescriptionPair{
					Code:        CodiceNonUtilizzabile,
					Description: "il codice inserito è scaduto e non più utilizzabile",
				},
			},
		},
	},
	TokenIdProviderType: &token.TokenIdProviderType{ProviderType: token.TokenIdProviderTypeExternal},
}
