package types

import (
	"filmogophery/internal/pkg/constant"
)

type (
	Token struct {
		AccessToken  string            `json:"accessToken"`
		RefreshToken string            `json:"refreshToken"`
		TokenType    string            `json:"tokenType"`
		ExpiresIn    int64             `json:"expiresIn"`
		ExpiresAt    constant.Datetime `json:"expiresAt"`
	}
)
