package gaefire

import (
	"github.com/eaglesakura/gaefire"
	fire_utils "github.com/eaglesakura/gaefire/utils"
	"testing"
	"github.com/stretchr/testify/assert"
)

func newTestServiceAccount() gaefire.FirebaseServiceAccount {
	fire := fire_utils.NewGaeFire()
	if json, err := fire.NewAssetManager().LoadFile("assets/firebase-admin.json"); err != nil {
		panic(err)
	} else {
		return fire.NewServiceAccount(json)
	}
}

func TestNewFirebaseServiceAccount(t *testing.T) {
	account := newTestServiceAccount()
	assert.NotNil(t, account)
}
