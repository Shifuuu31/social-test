package models

type Logininfo struct {
	Email    string `json:email`
	Password string `json:password`
}

type RegisterInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}
