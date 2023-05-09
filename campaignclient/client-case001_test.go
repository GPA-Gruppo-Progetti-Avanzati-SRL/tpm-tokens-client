package campaignclient_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-tokens-client/campaignclient"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-tokens-client/tokensclient"
)

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

// This file is being separated from the actual tests because in the current state of affairs it has to be aligned with the tokensazstore package one and the client one.
// easier to copy a file and change package instead of copying and pasting

var campaignTestCase001 = campaignclient.Campaign{
	Filters: campaignclient.Filters{
		Canale:   "APP",
		Servizio: "CC",
		Prodotto: "Start",
		Fase:     "Apertura",
	},
	CampaignType: campaignclient.Type{
		Code:        "MGM",
		Description: "member get member",
		Unique:      true,
		PromoCode:   "CPQ promo code",
		//TargetProducts: []campaignclient.ProductInfo{
		//	{
		//		Code:  "Codice prodotto interessato alla campagna (es. Start)",
		//		Ambit: "Codice servizio collegato (es. CC)",
		//	},
		//},
	},
	Title:       "L'amicizia ti premia",
	Description: "Porta un amico ed otterrete entrambi un conto gratuito per 12 mesi",
	AddInfo: campaignclient.AdditionalInfo{
		AltDescription:   "Porta un amico ed otterrete entrambi il seguente sconto",
		AwardDescription: "€0 / mese per 12 mesi",
	},
	Resources: []campaignclient.LinkedResource{
		{
			Type:        "regolemento",
			Name:        "Regolamento Campagna",
			ContentType: "application/pdf",
			Locations: []campaignclient.LinkedResourceLocation{
				{
					Type: "straight",
					Url:  "http://www.posteitaliane.it/campagne/BPMIFI/regolemento",
				},
			},
			Help: "Scarica il regolamento con tutte le informazioni della campagna",
		},
	},
	TokenContext: tokensclient.TokenContext{
		Id:            "BPMGM1",
		Pkey:          tokensclient.ContextPartitionKey,
		Platform:      "BP",
		Version:       tokensclient.TokenContextBaseVersion,
		Suspended:     false,
		BannerTokenId: "BPMIFI-BANNER",
		Timeline: tokensclient.Timeline{
			StartDate:      "20230101",
			EndDate:        "20230430",
			ExpirationMode: tokensclient.ExpirationModeDate,
		},
		StateMachine: tokensclient.StateMachine{
			States: []tokensclient.StateDefinition{
				{
					Code:        "[*]",
					StateType:   tokensclient.StartEndState,
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
					OutTransitions: []tokensclient.Transition{
						{
							To:          MGMStatusGenerato,
							Description: "il codice e' stato creato per il cf: {v:cf1}",
							Properties: []tokensclient.Property{
								{
									Name:           ParamNameCf,
									ValidationRule: "required",
									Help:           tokensclient.CodeDescriptionPair{Description: ParamCfHelpMessage},
								},
							},
							Bearers: []tokensclient.BearerRef{
								{
									Id:   BearerCF1ReferenceVariable,
									Role: "primary",
								},
							},
							ProcessVarDefinitions: []tokensclient.ProcessVarDefinition{
								{
									Name:        "cf1",
									Description: "codice fiscale utilizzatore",
									Value:       InputParamCf,
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
					StateType:   tokensclient.StateStd,
					Description: "Il codice è in attesa di un primo utilizzo da parte del codice fiscale dell'assegnatario",
					Actions: []tokensclient.ActionDefinition{
						{
							ActionId:   "in-action-verifica-15-cf",
							ActionType: tokensclient.ActionTypeIn,
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
					OutTransitions: []tokensclient.Transition{
						{
							To:          MGMStatusInAttesaAperturaPrimoConto,
							Description: "Il codice e' stato usato da {v:cf1} e il sistema e in attesa di apertura del primo conto",
							Properties: []tokensclient.Property{
								{
									Name:           ParamNameCf,
									ValidationRule: "required",
									Help:           tokensclient.CodeDescriptionPair{Description: ParamCfHelpMessage},
								},
								{
									Name:           ParamNameCanale,
									ValidationRule: "required",
									Help:           tokensclient.CodeDescriptionPair{Description: ParamCanaleHelpMessage},
								},
								{
									Name:           ParamNameProdotto,
									ValidationRule: "required",
									Help:           tokensclient.CodeDescriptionPair{Description: ParamProdottoHelpMessage},
								},
								{
									Name:           ParamNameFase,
									ValidationRule: "required",
									Help:           tokensclient.CodeDescriptionPair{Description: ParamFaseHelpMessage},
								},
								{
									Name:           ParamNameFunnelId,
									ValidationRule: "required",
									Help:           tokensclient.CodeDescriptionPair{Description: ParamFunnelIdHelpMessage},
								},
								{
									Name:           ParamNameNumeroPratica,
									ValidationRule: "required",
									Help:           tokensclient.CodeDescriptionPair{Description: ParamNumeroPraticaHelpMessage},
									Scope:          string(tokensclient.EventTypeNext),
								},
								{
									Name:           ParamNameServizio,
									ValidationRule: "required",
									Help:           tokensclient.CodeDescriptionPair{Description: ParamServizioHelpMessage},
									Scope:          string(tokensclient.EventTypeNext),
								},
								{
									Name:           ParamNameNumero,
									ValidationRule: "required",
									Help:           tokensclient.CodeDescriptionPair{Description: ParamNumeroHelpMessage},
									Scope:          string(tokensclient.EventTypeNext),
								},
							},
							Rules: []tokensclient.Rule{
								{
									Expression: "\"{$.cf}\" == \"{v:cf1}\"",
									Help:       tokensclient.CodeDescriptionPair{Description: "Il codice inserito non risulta associato al cliente {$.cf} ma al cliente {v:cf1}"},
								},
							},
							Bearers: []tokensclient.BearerRef{
								{
									Id:   BearerCF1ReferenceVariable,
									Role: "primary",
								},
							},
							ProcessVarDefinitions: []tokensclient.ProcessVarDefinition{
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
							TTL: tokensclient.TTLDefinition{
								Value: "1d",
							},
						},
					},
				},
				{
					Code:        MGMStatusInAttesaAperturaPrimoConto,
					StateType:   tokensclient.StateStd,
					Description: "il sistema e' in attesa del perfezionamento della pratica relativa al primo conto corrente",
					Help: tokensclient.CodeDescriptionPair{
						Code:        CodiceNonUtilizzabile,
						Description: "il codice inserito non e' ancora stato attivato",
					},
					OutTransitions: []tokensclient.Transition{
						{
							To:          MGMStatusPrimoContoAperto,
							Description: "il perfezionamento della pratica {v:numeroPratica1} di apertura del conto {v:servizio1}-{v:numero1} è terminata con successo",
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
							Bearers: []tokensclient.BearerRef{
								{
									Id:   BearerCF1ReferenceVariable,
									Role: "primary",
								},
							},
						},
						{
							To:          MGMStatusGenerato,
							Description: "il perfezionamento della pratica {v:numeroPratica1} di apertura del conto {v:servizio}-{v:numero} non e' andato a buon fine",
							Properties: []tokensclient.Property{
								{
									Name:           "result",
									ValidationRule: "required",
								},
							},
							Bearers: []tokensclient.BearerRef{
								{
									Id:   BearerCF1ReferenceVariable,
									Role: "primary",
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
					Code:        MGMStatusPrimoContoAperto,
					StateType:   tokensclient.StateStd,
					Description: "Il primo conto e' stato aperto con successo, il codice puo' essere ora utilizzato da un codice fiscale presentato dall'assegnatario del codice",
					Actions: []tokensclient.ActionDefinition{
						{
							ActionId:   "in-action-verifica-15-cf",
							ActionType: tokensclient.ActionTypeIn,
							Properties: map[string]interface{}{
								ParamNameCf: InputParamCf,
							},
						},
					},
					OutTransitions: []tokensclient.Transition{
						{
							To: MGMStatusInAttesaAperturaSecondoConto,
							Properties: []tokensclient.Property{
								{
									Name:           ParamNameCf,
									ValidationRule: "required",
									Help:           tokensclient.CodeDescriptionPair{Description: ParamCfHelpMessage},
								},
								{
									Name:           ParamNameCanale,
									ValidationRule: "required",
									Help:           tokensclient.CodeDescriptionPair{Description: ParamCanaleHelpMessage},
								},
								{
									Name:           ParamNameProdotto,
									ValidationRule: "required",
									Help:           tokensclient.CodeDescriptionPair{Description: ParamProdottoHelpMessage},
								},
								{
									Name:           ParamNameFase,
									ValidationRule: "required",
									Help:           tokensclient.CodeDescriptionPair{Description: ParamFaseHelpMessage},
								},
								{
									Name:           ParamNameFunnelId,
									ValidationRule: "required",
									Help:           tokensclient.CodeDescriptionPair{Description: ParamFunnelIdHelpMessage},
								},
								{
									Name:           ParamNameNumeroPratica,
									ValidationRule: "required",
									Help:           tokensclient.CodeDescriptionPair{Description: ParamNumeroPraticaHelpMessage},
									Scope:          string(tokensclient.EventTypeNext),
								},
								{
									Name:           ParamNameServizio,
									ValidationRule: "required",
									Help:           tokensclient.CodeDescriptionPair{Description: ParamServizioHelpMessage},
									Scope:          string(tokensclient.EventTypeNext),
								},
								{
									Name:           ParamNameNumero,
									ValidationRule: "required",
									Help:           tokensclient.CodeDescriptionPair{Description: ParamNumeroHelpMessage},
									Scope:          string(tokensclient.EventTypeNext),
								},
							},
							Bearers: []tokensclient.BearerRef{
								{
									Id:   BearerCF1ReferenceVariable,
									Role: "primary",
								},
								{
									Id:   BearerCF2ReferenceVariable,
									Role: "secondary",
								},
							},
							ProcessVarDefinitions: []tokensclient.ProcessVarDefinition{
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
							Rules: []tokensclient.Rule{
								{
									Expression: "\"{$.cf}\" != \"{v:cf1}\"",
									Help: tokensclient.CodeDescriptionPair{
										Description: "il codice fiscale dell'utilizzatore {$.cf} deve essere diverso dal codice fiscale dell'assegnatario del codice {v:cf1}",
									},
								},
							},
						},
					},
				},
				{
					Code:        MGMStatusInAttesaAperturaSecondoConto,
					StateType:   tokensclient.StateStd,
					Description: "il sistema e' in attesa del perfezionamento della pratica del secondo conto",
					Help: tokensclient.CodeDescriptionPair{
						Code:        CodiceNonUtilizzabile,
						Description: "il codice inserito e' stato gia' usato e in attesa di lavorazione",
					},
					OutTransitions: []tokensclient.Transition{
						{
							To: MGMStatusBruciato,
							Properties: []tokensclient.Property{
								{
									Name:           "result",
									ValidationRule: "required",
								},
							},
							Bearers: []tokensclient.BearerRef{
								{
									Id:   BearerCF1ReferenceVariable,
									Role: "primary",
								},
								{
									Id:   BearerCF2ReferenceVariable,
									Role: "secondary",
								},
							},
							Rules: []tokensclient.Rule{
								{
									Expression: "\"{$.result}\" == \"OK\"",
								},
							},
						},
						{
							To: MGMStatusPrimoContoAperto,
							Properties: []tokensclient.Property{
								{
									Name:           "result",
									ValidationRule: "required",
								},
							},
							Bearers: []tokensclient.BearerRef{
								{
									Id:   BearerCF1ReferenceVariable,
									Role: "primary",
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
					Code:        MGMStatusBruciato,
					StateType:   tokensclient.StateFinal,
					Description: "il codice e' stato utilizzato",
					Help: tokensclient.CodeDescriptionPair{
						Code:        CodiceNonUtilizzabile,
						Description: "il codice inserito e' gia' stato utilizzato",
					},

					BusinessView: tokensclient.CodeDescriptionPair{},
				},
				{
					Code:        MGMStatusExpired,
					StateType:   tokensclient.StateExpired,
					Description: "il codice inserito è scaduto e non più utilizzabile",
					Help: tokensclient.CodeDescriptionPair{
						Code:        CodiceNonUtilizzabile,
						Description: "il codice inserito è scaduto e non più utilizzabile",
					},
				},
			},
		},
		TokenIdProviderType: &tokensclient.TokenIdProviderType{ProviderType: tokensclient.TokenIdProviderTypeExternal},
	},
}
