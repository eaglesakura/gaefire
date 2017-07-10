package gaefire

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
	"time"
)

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
		kind: KindInfo{Name: "Default", Version: 1},
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
		log.Debugf(it.ctx, "Marshal error[%v]", key)
		return &DatastoreError{message: "Parse failed", errors: []error{err}}
	}

	item := &memcache.Item{
		Key:   key,
		Value: buf,
	}

	if it.expireDate != nil {
		item.Expiration = it.expireDate.Sub(time.Now())
	}

	return memcache.Set(it.ctx, item)
}

func (it *MemcacheLoadRequest) Load(result interface{}, createFunc func(result interface{}) error) error {
	key := it.getKey()
	item, err := memcache.Get(it.ctx, key)

	if err == nil {
		// デコードを行わせる
		return json.Unmarshal(item.Value, result)
	} else if err != nil && createFunc != nil {
		// memcahceに無いので、データを生成させる
		log.Debugf(it.ctx, "Memcache not found[%v]", key)
		createErr := createFunc(result)
		if createErr == nil {
			it.Save(result)
			return nil
		} else {
			return createErr
		}
	}

	return err
}
