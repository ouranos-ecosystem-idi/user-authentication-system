package entity

// SecureTokenAPI
// Summary: This is the struct which defines the SecureTokenAPI response entity.
type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
	UserId       string `json:"user_id"`
	ProjectId    string `json:"project_id"`
}
