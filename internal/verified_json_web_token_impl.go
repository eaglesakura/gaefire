package gaefire

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/eaglesakura/gaefire"
)

type VerifiedJsonWebTokenImpl struct {
	token *jwt.Token
}

func (it *VerifiedJsonWebTokenImpl) GetUserId() (string, error) {
	result, err := it.GetClaim("user_id")
	if err != nil {
		return "", err
	} else {
		return fmt.Sprintf("%v", result), nil
	}
}

/*
 * Firebaseユーザーを取得する
 */
func (it *VerifiedJsonWebTokenImpl) GetUser(result *gaefire.FirebaseUser) error {
	uid, err := it.GetUserId()
	if err != nil {
		return err
	}

	result.UniqueId = uid
	return nil
}

func (it *VerifiedJsonWebTokenImpl) GetAudience() (string, error) {
	result, err := it.GetClaim("aud")
	if err != nil {
		return "", err
	} else {
		return fmt.Sprintf("%v", result), nil
	}
}

func (it *VerifiedJsonWebTokenImpl) GetClaim(key string) (interface{}, error) {
	if result, ok := it.token.Claims.(jwt.MapClaims)[key]; !ok {
		return nil, errors.New(fmt.Sprintf("NotFound[%v]", key))
	} else {
		return result, nil
	}
}

func (it *VerifiedJsonWebTokenImpl) GetHeader(key string) (interface{}, error) {
	if result, ok := it.token.Header[key]; !ok {
		return nil, errors.New(fmt.Sprintf("NotFound[%v]", key))
	} else {
		return result, nil
	}
}
