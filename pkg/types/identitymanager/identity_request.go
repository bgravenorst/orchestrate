package identitymanager

type CreateIdentityRequest struct {
	Alias      string            `json:"alias" validate:"required" example:"personal-account"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

type ImportIdentityRequest struct {
	Alias      string            `json:"alias" validate:"required" example:"personal-account"`
	PrivateKey string            `json:"privateKey" validate:"required" example:"66232652FDFFD802B7252A456DBD8F3ECC0352BBDE76C23B40AFE8AEBD714E2D"`
	Attributes map[string]string `json:"attributes,omitempty"`
}