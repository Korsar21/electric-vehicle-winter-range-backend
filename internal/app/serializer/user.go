package serializer

type UserRegisterRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserAuthRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type UserAuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	UserID      uint   `json:"user_id"`
	Name        string `json:"name"`
	Role        string `json:"role"`
}
