package gaefire

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/url"
	"sync"

	"context"
	"sync/atomic"
)

type PublicKey struct {
	id        string
	publicKey *rsa.PublicKey
}

type PublicKeyGroup struct {
	keys map[string]*PublicKey // PublicKeys
}

/**
 * 公開鍵を管理する。
 *
 * 公開鍵は逐次変わるので、現状ではオンメモリにて管理される。
 */
type PublicKeystore struct {
	readTarget int32
	accountId  string // GCP Service account email.
	refreshUrl string // Publickey URL
	mutex      *sync.Mutex
	keystore   []*PublicKeyGroup
}

//
// create keystore
// @param refreshUrl "https://www.googleapis.com/robot/v1/metadata/x509/your-service-account@example.com"
//
func NewPublicKeystore(accountId string) *PublicKeystore {
	url := url.URL{}
	url.Path = "https://www.googleapis.com/robot/v1/metadata/x509/" + accountId
	return &PublicKeystore{
		readTarget: 0,
		accountId:  accountId,
		refreshUrl: url.EscapedPath(),
		mutex:      new(sync.Mutex),
		keystore:   []*PublicKeyGroup{nil, nil},
	}
}

//
// Find Public key
//
func (it *PublicKeystore) FindPublicKey(kid string) (*rsa.PublicKey, error) {
	read := it.keystore[atomic.LoadInt32(&it.readTarget)]
	if read != nil {
		value := read.keys[kid]
		if value != nil {
			return value.publicKey, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("Not found key[%v]", kid))
}

//
// Refresh Keystore
//
func (it *PublicKeystore) Refresh(ctx context.Context) error {

	// download public keys
	resp, err := newHttpClient(ctx).Get(it.refreshUrl)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		logError(fmt.Sprintf("CertRefresh failed err(%v) url(%v)\n", err.Error(), it.refreshUrl))
		return err
	}

	writeData := &PublicKeyGroup{
		keys: map[string]*PublicKey{},
	}

	// pull data & parse public key
	{
		var googlePublicKey interface{}
		buf, ioError := ioutil.ReadAll(resp.Body)
		if ioError != nil {
			logError(fmt.Sprintf("CertRefresh failed err(%v) url(%v)", ioError.Error(), it.refreshUrl))
			return err
		}

		ioError = json.Unmarshal(buf, &googlePublicKey)
		if ioError != nil {
			logError(fmt.Sprintf("CertRefresh failed err(%v) url(%v)", ioError.Error(), it.refreshUrl))
			return err
		}

		// to map
		for key, value := range googlePublicKey.(map[string]interface{}) {
			logDebug(fmt.Sprintf("Pull key id[%v] value[%v...]", key, value.(string)[:10]))
			parsedKey, keyError := jwt.ParseRSAPublicKeyFromPEM([]byte(value.(string)))
			if keyError != nil {
				logError(fmt.Sprintf("RSA Key failed err(%v) url(%v)", keyError.Error(), it.refreshUrl))
				return err
			}

			writeData.keys[key] = &PublicKey{
				id:        key,
				publicKey: parsedKey,
			}
		}
	}

	// swap read targets
	it.mutex.Lock()
	defer it.mutex.Unlock()

	writeTarget := (atomic.LoadInt32(&it.readTarget) + 1) % 2
	logInfo(fmt.Sprintf("Before Read Index[%v] Write Index[%v]", it.readTarget, writeTarget))
	it.keystore[writeTarget] = writeData
	atomic.StoreInt32(&it.readTarget, writeTarget)
	logInfo(fmt.Sprintf(" * Swap Read Index[%v]", it.readTarget))

	return nil
}
