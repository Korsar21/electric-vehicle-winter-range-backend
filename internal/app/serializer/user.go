package serializer

type UserRegisterRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
