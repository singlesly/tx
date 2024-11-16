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
	var val []byte
	err = item.Value(func(x []byte) error {
		val = append([]byte{}, x...)

		return nil
	})
	if err != nil {
		return nil, err
	}

	var txProto proto.Transaction
	p2p.DeserializeProtoMessage(val, &txProto)

	return &txProto, nil
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
