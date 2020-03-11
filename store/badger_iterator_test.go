package store_test

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"testing"
)

func TestBadgerStoreIterator(t *testing.T) {
	db, dir := newBadgerDB(t)
	defer os.RemoveAll(dir)
	defer db.Close()

	ks := make([][]byte, 30)
	vs := make([][]byte, len(ks))
	var i int

	prefix := []byte{0x61, 0x3a, 0x62, 0x3a}

	for i := 0; i < len(ks); i++ {
		s := strconv.Itoa(10 + i)
		k := append(prefix, s...)
		v := []byte(s)

		ks[i] = k
		vs[i] = v
		db.Set(k, v, 0)
	}

	iter := db.Iterate(prefix)
	defer iter.Done()

	for iter.Next() {
		item := iter.Item()
		assert.Equal(t, ks[i], item.Key())
		assert.Equal(t, vs[i], item.Value())
		assert.Zero(t, item.TTL())
		i++
	}
}
