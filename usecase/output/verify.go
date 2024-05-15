package output

// VerifyTokenResponse
// Summary: This is the structure which defines the verify token response.
type VerifyTokenResponse struct {
	OperatorID *string `json:"operatorId"`
}

// VerifyApiKeyResponse
// Summary: This is the structure which defines the verify API key response.
type VerifyApiKeyResponse struct {
	IsAPIKeyValid    bool `json:"isApiKeyValid"`
	IsIPAddressValid bool `json:"isIpAddressValid"`
}
