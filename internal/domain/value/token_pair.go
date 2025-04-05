package value

type TokenPair struct {
	AccessToken  `json:"accessToken"`
	RefreshToken `json:"refreshToken"`
}
