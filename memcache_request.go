package gaefire

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"time"
)

var memcache = cache.New(60*time.Minute, 10*time.Minute)

/**
 * Memcache/Datastoreに保存するデータのKind情報を定義する
 */
type KindInfo struct {
	/**
	 * Kind名, Keyの一部に利用される
	 */
	Name string

	/**
	 * バージョン番号, Keyの一部に利用される
	 */
	Version int
}

type MemcacheLoadRequest struct {
	ctx        context.Context
	kind       KindInfo
	uniqueId   string
	expireDate *time.Time
}

func NewMemcacheRequest(ctx context.Context) *MemcacheLoadRequest {
	return &MemcacheLoadRequest{
		ctx:  ctx,
		kind: KindInfo{Name: "Default", Version: 2},
	}
}

func (it *MemcacheLoadRequest) SetExpireDate(date time.Time) *MemcacheLoadRequest {
	it.expireDate = &date
	return it
}

func (it *MemcacheLoadRequest) SetKindInfo(kind KindInfo) *MemcacheLoadRequest {
	it.kind = kind
	return it
}

func (it *MemcacheLoadRequest) SetId(id string) *MemcacheLoadRequest {
	it.uniqueId = id
	return it
}

func (it *MemcacheLoadRequest) getKey() string {
	return fmt.Sprintf("internal.gaefire.%v:%v.%v", it.kind.Name, it.kind.Version, it.uniqueId)
}

func (it *MemcacheLoadRequest) Save(value interface{}) error {
	key := it.getKey()
	buf, err := json.Marshal(value)
	if err != nil {
		logDebug(fmt.Sprintf("Marshal error[%v]", key))
		return &DatastoreError{message: "Parse failed", errors: []error{err}}
	}

	memcache.Set(key, buf, it.expireDate.Sub(time.Now()))
	return nil
}

func (it *MemcacheLoadRequest) Load(result interface{}, createFunc func(result interface{}) error) error {
	key := it.getKey()
	if buf, found := memcache.Get(key); found {
		// found memcache.
		return json.Unmarshal(buf.([]byte), result)
	} else {
		// not found memcache.
		logDebug(fmt.Sprintf("Memcache not found[%v]", key))
		createErr := createFunc(result)
		if createErr == nil {
			return it.Save(result)
		} else {
			return createErr
		}
	}
}
