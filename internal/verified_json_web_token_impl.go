package internal

import (
	"github.com/dgrijalva/jwt-go"
	"fmt"
)

type VerifiedJsonWebTokenImpl struct {
	token *jwt.Token
}

func (it *VerifiedJsonWebTokenImpl)GetUserId() (string, error) {
	result, err := it.GetClaim("user_id")
	if err != nil {
		return "", err
	} else {
		return fmt.Sprintf("%v", result), nil
	}
}

func (it *VerifiedJsonWebTokenImpl)GetProjectId() (string, error) {
	result, err := it.GetClaim("aud")
	if err != nil {
		return "", err
	} else {
		return result, nil
	}
}

func (it *VerifiedJsonWebTokenImpl)GetClaim(key string) (interface{}, error) {
	result, err := it.token.Claims.(jwt.MapClaims)[key]
	if err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (it *VerifiedJsonWebTokenImpl)GetHeader(key string, def string) string {
	result, ok := it.token.Header[key]
	if !ok {
		return def
	} else {
		return result.(string)
	}
}