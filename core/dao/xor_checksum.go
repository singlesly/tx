package dao

import (
	"fmt"
	"github.com/dgraph-io/badger"
)

func GetXorChecksum(db *badger.DB, key string) (string, error) {
	result := make([]byte, 0)

	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return fmt.Errorf("[GetXorChecksum] %w", err)
		}

		item.Value(func(val []byte) error {
			result = append(result, val...)

			return nil
		})
		return nil
	})

	return string(result), err
}

func InsertXorChecksum(db *badger.DB, key string, value string) error {
	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), []byte(value))
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
