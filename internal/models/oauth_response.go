package models

type OAuthResponse struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string   `json:"name"`
	Picture       string   `json:"picture"`
}
