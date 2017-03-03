package gaefire

import (
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"github.com/eaglesakura/gaefire"
)

/**
 * Token生成用の構造体
 */
type TokenSourceModel struct {
	jwt.StandardClaims
	Uid    string`json:"uid,omitempty"`
	Scope  string`json:"scope,omitempty"`
	Claims map[string]string`json:"claims,omitempty"`
}

/**
 * Json Web Tokenを生成するGenerator
 */
type JsonWebTokenGeneratorImpl struct {
	service   gaefire.FirebaseServiceAccount
	source    TokenSourceModel
	headers   map[string]string
	lastError error;
}

/**
 * creamを登録する
 */
func (it *JsonWebTokenGeneratorImpl)AddClaim(key string, value interface{}) gaefire.JsonWebTokenGenerator {
	if it.lastError == nil {
		it.source.Claims[key] = fmt.Sprintf("%v", value)
	}
	return it
}

/**
 * JWTのヘッダに情報を付与する
 *
 * ex) AddClaim("kid", "your.private.key.id");
 */
func (it *JsonWebTokenGeneratorImpl)AddHeader(key string, value interface{}) gaefire.JsonWebTokenGenerator {
	if it.lastError != nil {
		return it
	}

	it.headers[key] = fmt.Sprintf("%v", value)
	return it
}

/**
 * Json Web Tokenの仕様に従って署名された値を生成する。
 *
 * `base64(header).base64(token).base64(sign)`
 */
func (it *JsonWebTokenGeneratorImpl)Generate() (string, error) {
	if it.lastError != nil {
		return "", it.lastError
	}

	// Gen Token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, it.source)

	for key, value := range it.headers {
		jwtToken.Header[key] = value;
	}

	// Gen JWT
	signed, err := jwtToken.SignedString(it.service.GetPrivateKey())
	if err != nil {
		it.lastError = err
		return "", it.lastError
	}

	return signed, nil

}
