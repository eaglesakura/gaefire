package internal

import (
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"errors"
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
		return fmt.Sprintf("%v", result), nil
	}
}

func (it *VerifiedJsonWebTokenImpl)GetClaim(key string) (interface{}, error) {
	if result, ok := it.token.Claims.(jwt.MapClaims)[key]; !ok {
		return nil, errors.New(fmt.Sprintf("NotFound[%v]", key))
	} else {
		return result, nil
	}
}

func (it *VerifiedJsonWebTokenImpl)GetHeader(key string) (interface{}, error) {
	if result, ok := it.token.Header[key]; !ok {
		return nil, errors.New(fmt.Sprintf("NotFound[%v]", key))
	} else {
		return result, nil
	}
}