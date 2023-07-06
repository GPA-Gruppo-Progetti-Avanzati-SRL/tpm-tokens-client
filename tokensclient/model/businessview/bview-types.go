package businessview

type BearerRef struct {
	Id   string `yaml:"id,omitempty" mapstructure:"id,omitempty" json:"id,omitempty"`
	Role string `yaml:"role,omitempty" mapstructure:"role,omitempty" json:"role,omitempty"`
}

type EventState struct {
	Code        string `yaml:"code,omitempty" mapstructure:"code,omitempty" json:"code,omitempty"`
	Description string `yaml:"description,omitempty" mapstructure:"description,omitempty" json:"description,omitempty"`
	Pending     bool   `yaml:"pending,omitempty" mapstructure:"pending,omitempty" json:"pending,omitempty"`
}

type Event struct {
	Description string     `yaml:"description,omitempty" mapstructure:"description,omitempty" json:"description,omitempty"`
	Ts          string     `yaml:"ts,omitempty" mapstructure:"ts,omitempty" json:"ts,omitempty"`
	Name        string     `yaml:"name,omitempty" mapstructure:"name,omitempty" json:"name,omitempty"`
	State       EventState `yaml:"state,omitempty" mapstructure:"state,omitempty" json:"state,omitempty"`
}

type Property struct {
	Name        string      `yaml:"name,omitempty" mapstructure:"name,omitempty" json:"name,omitempty"`
	Value       interface{} `yaml:"value,omitempty" mapstructure:"value,omitempty" json:"value,omitempty"`
	Description string      `yaml:"description,omitempty" mapstructure:"description,omitempty" json:"description,omitempty"`
}

type Token struct {
	Id         string      `yaml:"id,omitempty" mapstructure:"id,omitempty" json:"id,omitempty"`
	ContextId  string      `yaml:"ctx-id,omitempty" mapstructure:"ctx-id,omitempty" json:"ctx-id,omitempty"`
	Typ        string      `yaml:"type,omitempty" mapstructure:"type,omitempty" json:"type,omitempty"`
	ExpiryTs   string      `yaml:"expiry-ts,omitempty" mapstructure:"expiry-ts,omitempty" json:"expiry-ts,omitempty"`
	Events     []Event     `yaml:"events,omitempty" mapstructure:"events,omitempty" json:"events,omitempty"`
	Properties []Property  `yaml:"properties,omitempty" mapstructure:"properties,omitempty" json:"properties,omitempty"`
	Bearers    []BearerRef `yaml:"bearers,omitempty" mapstructure:"bearers,omitempty" json:"bearers,omitempty"`
}

type TokenRef struct {
	Token Token  `yaml:"token,omitempty" mapstructure:"token,omitempty" json:"token,omitempty"`
	Role  string `yaml:"role,omitempty" mapstructure:"role,omitempty" json:"role,omitempty"`
}

type Bearer struct {
	Id         string     `yaml:"id,omitempty" mapstructure:"id,omitempty" json:"id,omitempty"`
	ContextId  string     `yaml:"ctx-id,omitempty" mapstructure:"ctx-id,omitempty" json:"ctx-id,omitempty"`
	TokenRefs  []TokenRef `yaml:"tok-refs,omitempty" mapstructure:"tok-refs,omitempty" json:"tok-refs,omitempty"`
	Properties []Property `yaml:"properties,omitempty" mapstructure:"properties,omitempty" json:"properties,omitempty"`
}

type Actor struct {
	ActorId string   `yaml:"actor-id,omitempty" mapstructure:"actor-id,omitempty" json:"actor-id,omitempty"`
	Bearers []Bearer `yaml:"contexts,omitempty" mapstructure:"contexts,omitempty" json:"contexts,omitempty"`
}
