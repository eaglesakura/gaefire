package gaefire

import (
	"crypto/rsa"
	"sync"
	"io/ioutil"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"errors"
	"fmt"
	"net/url"

	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
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
	readTarget int
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
		readTarget:0,
		accountId:accountId,
		refreshUrl:url.EscapedPath(),
		mutex : new(sync.Mutex),
		keystore:[]*PublicKeyGroup{nil, nil},
	}
}

//
// Find Public key
//
func (it *PublicKeystore)FindPublicKey(kid string) (*rsa.PublicKey, error) {
	read := it.keystore[it.readTarget]
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
func (it *PublicKeystore)Refresh(ctx context.Context) error {

	// download public keys
	resp, err := newHttpClient(ctx).Get(it.refreshUrl)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		log.Errorf(ctx, "CertRefresh failed err(%v) url(%v)\n", err.Error(), it.refreshUrl)
		return err;
	}

	writeData := &PublicKeyGroup{
		keys:map[string]*PublicKey{},
	}

	// pull data & parse public key
	{
		var googlePublicKey  interface{};
		buf, ioEror := ioutil.ReadAll(resp.Body)
		if ioEror != nil {
			log.Errorf(ctx, "CertRefresh failed err(%v) url(%v)", ioEror.Error(), it.refreshUrl)
			return err
		}

		ioEror = json.Unmarshal(buf, &googlePublicKey);
		if ioEror != nil {
			log.Errorf(ctx, "CertRefresh failed err(%v) url(%v)", ioEror.Error(), it.refreshUrl)
			return err
		}

		// to map
		for key, value := range googlePublicKey.(map[string]interface{}) {
			log.Debugf(ctx, "Pull key id[%v] value[%v...]", key, value.(string)[:10])
			parsedKey, keyError := jwt.ParseRSAPublicKeyFromPEM([]byte(value.(string)))
			if keyError != nil {
				log.Errorf(ctx, "RSA Key failed err(%v) url(%v)", keyError.Error(), it.refreshUrl)
				return err
			}

			writeData.keys[key] = &PublicKey{
				id:key,
				publicKey:parsedKey,
			}
		}
	}

	// swap read targets
	it.mutex.Lock()
	defer it.mutex.Unlock()

	writeTarget := (it.readTarget + 1) % 2

	it.keystore[writeTarget] = writeData
	it.readTarget = writeTarget

	return nil
}
