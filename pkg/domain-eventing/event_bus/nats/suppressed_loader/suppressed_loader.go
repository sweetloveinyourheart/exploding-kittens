package suppressedloader

import (
	"fmt"

	"golang.org/x/sync/singleflight"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/ttlcache"
)

// SuppressedLoader wraps another Loader and suppresses duplicate
// calls to its Load method.
type SuppressedLoader[K comparable, V any] struct {
	loader ttlcache.Loader[K, V]
	group  *singleflight.Group
}

// NewSuppressedLoader creates a new instance of suppressed loader.
// If the group parameter is nil, a newly created instance of
// *singleflight.Group is used.
func NewSuppressedLoader[K comparable, V any](loader ttlcache.Loader[K, V], group *singleflight.Group) *SuppressedLoader[K, V] {
	if group == nil {
		group = &singleflight.Group{}
	}

	return &SuppressedLoader[K, V]{
		loader: loader,
		group:  group,
	}
}

// Load executes a custom item retrieval logic and returns the item that
// is associated with the key.
// It returns nil if the item is not found/valid.
// It also ensures that only one execution of the wrapped Loader's Load
// method is in-flight for a given key at a time.
func (l *SuppressedLoader[K, V]) Load(c *ttlcache.Cache[K, V], key K) *ttlcache.Item[K, V] {
	// there should be a better/generic way to create a
	// singleflight Group's key. It's possible that a generic
	// singleflight.Group will be introduced with/in go1.19+
	strKey := fmt.Sprint(key)

	// the error can be discarded since the singleflight.Group
	// itself does not return any of its errors, it returns
	// the error that we return ourselves in the func below, which
	// is also nil
	res, _, _ := l.group.Do(strKey, func() (interface{}, error) {
		item := l.loader.Load(c, key)
		if item == nil {
			return nil, nil //nolint:nilnil
		}

		return item, nil
	})
	if res == nil {
		return nil
	}

	return res.(*ttlcache.Item[K, V])
}
