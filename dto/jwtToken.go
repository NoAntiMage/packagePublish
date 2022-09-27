package dto

type LoginToken struct {
	User       string `json:"User"`
	Area       string `json:"Area"`
	LoginToken string `json:"LoginToken"`
}

type JwtToken struct {
	User     string `json:"User"`
	JwtToken string `json:"JwtToken"`
}

type RpcTokenUpdate struct {
	User       string `json:"User" validate:"required"`
	ExpireTime int    `json:"ExpireTime"`
}
