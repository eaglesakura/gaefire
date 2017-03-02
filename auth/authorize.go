package gaefire

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"errors"
	"fmt"
	"strings"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

//type FirebaseServiceAccount struct {
//	serviceAccount internal.ServiceAccountJson
//
//	// service account private key(for signing)
//	privateKey     *rsa.PrivateKey
//	publicKeystore *PublicKeystore
//}


// load service.json
//func NewServiceFromFile(path string) (*FirebaseServiceAccount, error) {
//	buf, err := ioutil.ReadFile(path);
//	if err != nil {
//		return nil, err
//	}
//
//	account := internal.ServiceAccountJson{}
//	err = json.Unmarshal(buf, &account)
//	if err != nil {
//		return nil, err;
//	}
//
//	return NewServiceAccount(&account)
//}

//func NewServiceAccount(serviceAccount *internal.ServiceAccountJson) (*FirebaseServiceAccount, error) {
//	parsedPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(serviceAccount.PrivateKey));
//	if err != nil || parsedPrivateKey == nil {
//		return nil, err;
//	}
//
//	result := &FirebaseServiceAccount{
//		serviceAccount:*serviceAccount,
//		privateKey:parsedPrivateKey,
//		publicKeystore:NewKeystore(serviceAccount.ClientEmail),
//	}
//
//	return result, nil
//}
