package types

type Base struct {
	ProviderType  string `json:"_provider" diff:"-"`
	ComponentType string `json:"_component" diff:"-"`
	ComponentID   string `json:"_component_id" diff:"_component_id,immutable"`
	State         string `json:"_state" diff:"-"`
	Action        string `json:"_action" diff:"-"`
	Service       string `json:"service" diff:"-"`
	ErrorMessage  string `json:"error,omitempty" diff:"-"`
	Subject       string `json:"-" diff:"-"`
	Body          []byte `json:"-" diff:"-"`
	CryptoKey     string `json:"-" diff:"-"`
}
