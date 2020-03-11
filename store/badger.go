package store

import (
	"github.com/dgraph-io/badger"
	"time"
)

type BadgerStore struct {
	db   *badger.DB
	quit chan struct{}
}

func (b *BadgerStore) Size() int64 {
	lsm, vlog := b.db.Size()
	return lsm + vlog
}

func (b *BadgerStore) Set(key, val []byte, ttl time.Duration) {
	txn := b.db.NewTransaction(true)
	defer txn.Commit()

	entry := badger.NewEntry(key, val)
	if ttl > 0 {
		entry.WithTTL(ttl)
	}
	_ = txn.SetEntry(entry)
}

func (b *BadgerStore) Get(key []byte) ([]byte, bool) {
	txn := b.db.NewTransaction(false)
	item, err := txn.Get(key)
	if err != nil {
		return nil, false
	}
	data, err := item.ValueCopy(nil)
	if err != nil {
		return nil, false
	}
	return data, true
}

func (b *BadgerStore) Delete(key []byte) {
	txn := b.db.NewTransaction(true)
	defer txn.Commit()
	_ = txn.Delete(key)
}

func (b *BadgerStore) Iterate(key []byte) Iterator {
	txn := b.db.NewTransaction(false)

	opts := badger.IteratorOptions{}
	opts.PrefetchSize = 10
	opts.Prefix = key

	iter := txn.NewIterator(opts)

	bi := &BadgerIterator{iter: iter}
	bi.Seek(nil)
	return bi
}

func (b *BadgerStore) Close() {
	close(b.quit)
	b.db.Close()
}

func (b *BadgerStore) Init() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				b.db.RunValueLogGC(0.7)
			case <-b.quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func NewBadgerStore(path string) (Store, error) {
	opt := badger.DefaultOptions(path)
	//opt.Logger = &badger.Logger()

	db, err := badger.Open(opt)
	if err != nil {
		return nil, err
	}
	quit := make(chan struct{}, 1)

	store := &BadgerStore{db, quit}
	store.Init()
	return store, nil
}
