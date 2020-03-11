package store

import "time"

type Store interface {
	Init()
	Size() int64
	Set(key, val []byte, ttl time.Duration)
	Get(key []byte) (data []byte, found bool)
	Delete(key []byte)
	Iterate(key []byte) Iterator
	Close()
}

type Iterator interface {
	Seek(key []byte)
	Next() bool
	Item() Item
	Done()
}

type Item interface {
	Key() []byte
	Value() []byte
	TTL() time.Time
}
