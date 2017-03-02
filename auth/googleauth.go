package gaefire

import (
	"encoding/json"
	"io/ioutil"
	"github.com/dgrijalva/jwt-go"
	"errors"
	"fmt"
	"google.golang.org/appengine/urlfetch"
	"golang.org/x/net/context"
)

const (
	_GOOGLE_TOKEN_VERIFY_URL = "https://www.googleapis.com/oauth2/v3/tokeninfo?id_token="
)

type GoogleVerifyOption struct {
	TrustedClients []string // ClientID... hoge-fuga-uid.apps.googleusercontent.com
}

func (id *FirebaseServiceAccount)VerifyGoogleLogin(ctx context.Context, googleLoginToken string, opt *GoogleVerifyOption) (*VerifiedToken, *VerifyError) {
	client := urlfetch.Client(ctx)
	resp, err := client.Get(_GOOGLE_TOKEN_VERIFY_URL + googleLoginToken)
	defer resp.Body.Close()

	if err != nil {
		return nil, newTokenError(err)
	}

	buffer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, newTokenError(err)
	}
	var mapping jwt.MapClaims;
	err = json.Unmarshal(buffer, &mapping)
	if err != nil {
		return nil, newTokenError(err)
	}

	result := &VerifiedToken{
		token:&jwt.Token{
			Claims:mapping,
		},
	}

	projectId, err := result.GetProjectId()
	if err != nil {
		return nil, newTokenError(err)
	}

	trust := false
	for _, prj := range opt.TrustedClients {
		if projectId == prj {
			trust = true
		}
	}

	if !trust {
		return nil, newTokenError(errors.New(fmt.Sprintf("Token not trusted, ClientId[%v]", projectId)))
	}

	return result, nil
}
