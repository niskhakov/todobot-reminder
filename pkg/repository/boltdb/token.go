package boltdb

import (
	"encoding/binary"
	"errors"

	"github.com/boltdb/bolt"
	"github.com/niskhakov/todobot-reminder/pkg/repository"
)

type TokenRepository struct {
	db *bolt.DB
}

func NewTokenRepository(db *bolt.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Save(chatID int64, token string, bucket repository.Bucket) error {
	err := r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put(intToBytes(chatID), []byte(token))
	})

	return err
}

func (r *TokenRepository) Get(chatID int64, bucket repository.Bucket) (string, error) {
	var token string
	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		data := b.Get(intToBytes(chatID))
		token = string(data)
		return nil
	})

	if err != nil {
		return "", err
	}

	if token == "" {
		return "", errors.New("token not found")
	}

	return token, nil
}

func (r *TokenRepository) ForEach(bucket repository.Bucket, fnc repository.IterateFunc, accumulator interface{}) error {
	return r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			fnc(bytesToInt(k), string(v), accumulator)
		}

		return nil
	})
}

func intToBytes(v int64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(v))
	return buf
}

func bytesToInt(b []byte) int64 {
	return int64(binary.BigEndian.Uint64(b))
}
