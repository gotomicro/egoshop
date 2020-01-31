package trans

type ReqOauthLogin struct {
	Name string `json:"userName" binding:"required"`
	Pwd  string `json:"password" binding:"required"`
	Type string `json:"type"`
}
