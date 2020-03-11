package store_test

import (
	"github.com/Michael2008S/etherpad4go/store"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func newBadgerDB(t *testing.T) (store.Store, string) {
	dir, err := ioutil.TempDir("", "badger")
	assert.Nil(t, err)

	tmp := filepath.Join(dir, "db")
	d, err := store.NewBadgerStore(tmp)

	assert.Nil(t, err)
	return d, dir
}

func TestBadgerStore_Close(t *testing.T) {

}

func TestBadgerStore_Delete(t *testing.T) {
	db, dir := newBadgerDB(t)
	defer os.RemoveAll(dir)
	defer db.Close()

	o := []byte{0x61}

	b, ok := db.Get(o)
	assert.False(t, ok)
	assert.Nil(t, b)

	db.Set(o, o, 0)

	b, ok = db.Get(o)
	assert.True(t, ok)
	assert.Equal(t, o, b)

	db.Delete(o)

	b, ok = db.Get(o)
	assert.False(t, ok)
	assert.Nil(t, b)
}

func TestBadgerStore_Get(t *testing.T) {
	db, dir := newBadgerDB(t)
	defer os.RemoveAll(dir)
	defer db.Close()

	arr := make([][]byte, 500)

	for i := 0; i < len(arr); i++ {
		s := strconv.Itoa(i)
		arr[i] = []byte(s)
		db.Set(arr[i], arr[i], 0)
	}
	for i := 0; i < len(arr); i++ {
		o, found := db.Get(arr[i])
		assert.True(t, found)
		assert.Equal(t, arr[i], o)
	}
}


func TestBadgerStore_Iterate(t *testing.T) {

}

func TestBadgerStore_Set(t *testing.T) {
	db, dir := newBadgerDB(t)
	defer os.RemoveAll(dir)
	defer db.Close()

	o := []byte{0x62}
	db.Set(o, o, time.Second*2)

	b, ok := db.Get(o)
	assert.True(t, ok)
	assert.Equal(t, o, b)

	time.Sleep(time.Second*2 + time.Millisecond)

	b, ok = db.Get(o)
	assert.False(t, ok)
	assert.Nil(t, b)
}

func TestBadgerStore_Size(t *testing.T) {

}

func TestNewBadgerStore(t *testing.T) {
	db, dir := newBadgerDB(t)
	defer os.RemoveAll(dir)
	defer db.Close()

	assert.NotEmpty(t, dir)
	assert.NotNil(t, db)
}
