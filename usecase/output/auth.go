package output

import "strings"

// LoginResponse
// Summary: This is the structure which defines the login response.
type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// Mask
// Summary: This is the function which masks the confidential information.
func (o *LoginResponse) Mask() {
	o.RefreshToken = strings.Repeat("*", len(o.RefreshToken))
}

// RefreshResponse
// Summary: This is the structure which defines the refresh response.
type RefreshResponse struct {
	AccessToken string `json:"accessToken"`
}
