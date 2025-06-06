package ttlcache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_newItem(t *testing.T) {
	item := newItem("key", 123, time.Hour, false)
	require.NotNil(t, item)
	assert.Equal(t, "key", item.key)
	assert.Equal(t, 123, item.value)
	assert.Equal(t, time.Hour, item.ttl)
	assert.Equal(t, int64(-1), item.version)
	assert.WithinDuration(t, time.Now().Add(time.Hour), item.expiresAt, time.Minute)
}

func Test_Item_update(t *testing.T) {
	item := Item[string, string]{
		expiresAt: time.Now().Add(-time.Hour),
		value:     "hello",
		version:   0,
	}

	item.update("test", -1, time.Hour)
	assert.Equal(t, "test", item.value)
	assert.Equal(t, time.Hour, item.ttl)
	assert.Equal(t, int64(1), item.version)
	assert.WithinDuration(t, time.Now().Add(time.Hour), item.expiresAt, time.Minute)

	item.update("hi", -1, NoTTL)
	assert.Equal(t, "hi", item.value)
	assert.Equal(t, NoTTL, item.ttl)
	assert.Equal(t, int64(2), item.version)
	assert.Zero(t, item.expiresAt)
}

func Test_Item_touch(t *testing.T) {
	var item Item[string, string]
	item.touch()
	assert.Equal(t, int64(0), item.version)
	assert.Zero(t, item.expiresAt)

	item.ttl = time.Hour
	item.touch()
	assert.Equal(t, int64(0), item.version)
	assert.WithinDuration(t, time.Now().Add(time.Hour), item.expiresAt, time.Minute)
}

func Test_Item_IsExpired(t *testing.T) {
	// no ttl
	item := Item[string, string]{
		expiresAt: time.Now().Add(-time.Hour),
	}

	assert.False(t, item.IsExpired())

	// expired
	item.ttl = time.Hour
	assert.True(t, item.IsExpired())

	// not expired
	item.expiresAt = time.Now().Add(time.Hour)
	assert.False(t, item.IsExpired())
}

func Test_Item_Key(t *testing.T) {
	item := Item[string, string]{
		key: "test",
	}

	assert.Equal(t, "test", item.Key())
}

func Test_Item_Value(t *testing.T) {
	item := Item[string, string]{
		value: "test",
	}

	assert.Equal(t, "test", item.Value())
}

func Test_Item_TTL(t *testing.T) {
	item := Item[string, string]{
		ttl: time.Hour,
	}

	assert.Equal(t, time.Hour, item.TTL())
}

func Test_Item_ExpiresAt(t *testing.T) {
	now := time.Now()
	item := Item[string, string]{
		expiresAt: now,
	}

	assert.Equal(t, now, item.ExpiresAt())
}

func Test_Item_Version(t *testing.T) {
	item := Item[string, string]{version: 5}
	assert.Equal(t, int64(5), item.Version())
}
