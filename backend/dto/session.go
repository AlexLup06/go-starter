package dto

import "encoding/json"

type TokenCookie struct {
	SessionId string `json:"session_id"`
	Token     string `json:"token"`
}

func (s *TokenCookie) ToJsonString() (string, error) {
	bytes, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func NewTokenCookie(jsonStr string) (TokenCookie, error) {
	var result TokenCookie
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return TokenCookie{}, err
	}
	return result, nil
}
