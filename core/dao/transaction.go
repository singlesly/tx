package dao

import (
	"errors"
	"github.com/dgraph-io/badger"
	"transaction/p2p"
	"transaction/proto"
)

const TransactionTablePrefix = "transactions"

func GetTransactionKey(txId string) []byte {
	return []byte(TransactionTablePrefix + ":" + txId)
}

func HasTransactions(db *badger.DB) bool {
	dbTx := db.NewTransaction(false)

	opts := badger.DefaultIteratorOptions
	opts.Prefix = []byte(TransactionTablePrefix)
	iterator := dbTx.NewIterator(opts)
	defer iterator.Close()

	iterator.Rewind()
	return iterator.Valid()
}

func GetTransaction(db *badger.DB, txId string) (*proto.Transaction, error) {
	dbTx := db.NewTransaction(false)
	item, err := dbTx.Get(GetTransactionKey(txId))
	if err != nil {
		return nil, err
	}

	txProto := ItemToProtoTransaction(item)

	return txProto, nil
}

func GetTransactions(db *badger.DB) ([]*proto.Transaction, error) {
	result := make([]*proto.Transaction, 0)

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte(TransactionTablePrefix)
		opts.PrefetchSize = 1000

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			result = append(result, ItemToProtoTransaction(it.Item()))
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, err
}

func ExistsTransaction(db *badger.DB, txId string) bool {
	dbTx := db.NewTransaction(false)
	_, err := dbTx.Get(GetTransactionKey(txId))
	if errors.Is(err, badger.ErrKeyNotFound) {
		return false
	}

	return true
}

func InsertTransaction(db *badger.DB, tx *proto.Transaction) {
	dbTx := db.NewTransaction(true)
	err := dbTx.Set(GetTransactionKey(tx.TxId), p2p.SerializeProtoMessage(tx))
	if err != nil {
		return
	}
	err = dbTx.Commit()
	if err != nil {
		return
	}
}

func GetLastNTransactions(db *badger.DB, count int) []*proto.Transaction {
	result := make([]*proto.Transaction, 0)

	_ = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Reverse = true
		opts.PrefetchSize = count
		opts.Prefix = []byte(TransactionTablePrefix)

		it := txn.NewIterator(opts)

		for it.Rewind(); it.Valid(); it.Next() {
			result = append(result, ItemToProtoTransaction(it.Item()))
		}
		return nil
	})

	return result
}

func ItemToProtoTransaction(item *badger.Item) *proto.Transaction {
	var transaction proto.Transaction

	err := item.Value(func(val []byte) error {
		p2p.DeserializeProtoMessage(val, &transaction)

		return nil
	})
	if err != nil {
		return nil
	}

	return &transaction
}
