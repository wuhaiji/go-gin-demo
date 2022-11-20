package entity

type UserInfo struct {
	Password string `json:"password" form:"password"`
	Username string `json:"username" form:"username"`
}
