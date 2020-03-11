package store

import (
	"github.com/dgraph-io/badger"
	"time"
)

type BadgerIterator struct {
	iter *badger.Iterator
	curr int
}

func (b *BadgerIterator) Seek(key []byte) {
	b.curr = -1
	b.iter.Seek(key)
}

func (b *BadgerIterator) Next() bool {
	if b.curr > -1 {
		b.iter.Next()
	}
	b.curr++
	return b.iter.Valid()
}

func (b *BadgerIterator) Item() Item {
	return &BadgerIteratorItem{b.iter.Item()}
}

func (b *BadgerIterator) Done() {
	b.iter.Close()
}

type BadgerIteratorItem struct {
	item *badger.Item
}

func (b *BadgerIteratorItem) Key() []byte {
	return b.item.Key()
}

func (b *BadgerIteratorItem) Value() []byte {
	t, err := b.item.ValueCopy(nil)
	if err != nil {
		return nil
	}
	return t
}

func (b *BadgerIteratorItem) TTL() time.Time {
	ex := b.item.ExpiresAt()
	if ex > 0 {
		sec := int64(b.item.ExpiresAt())
		return time.Unix(sec, 0)
	}
	return time.Time{}
}
