package structs

type UserInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}
