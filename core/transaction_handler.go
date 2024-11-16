package core

import (
	"github.com/dgraph-io/badger"
	proto2 "google.golang.org/protobuf/proto"
	"log"
	"transaction/core/dao"
	"transaction/p2p"
	"transaction/proto"
)

type TransactionHandler struct {
	Db *badger.DB
}

func NewTransactionHandler(peer p2p.NodePeer, db *badger.DB) *TransactionHandler {

	handler := &TransactionHandler{
		Db: db,
	}

	peer.Subscribe("transactions", func(message []byte) {
		var transaction proto.Transaction
		proto2.Unmarshal(message, &transaction)

		handler.handleNewTransaction(&transaction)
	})

	return handler
}

func (h *TransactionHandler) handleNewTransaction(transaction *proto.Transaction) {

	if dao.ExistsTransaction(h.Db, transaction.TxId) {
		log.Printf("transaction %s skipped", transaction.TxId)
	} else if len(transaction.TxRefs) == 0 && !dao.HasTransactions(h.Db) {
		log.Printf("transaction %s applied", transaction.TxId)
		dao.InsertTransaction(h.Db, transaction)
		return
	} else if dao.HasTransactions(h.Db) && len(transaction.TxRefs) != 0 {
		for _, txRef := range transaction.TxRefs {
			if !dao.ExistsTransaction(h.Db, txRef) {
				log.Printf("tx ref %s not found. do not apply transaction", txRef)

				return
			}
		}

		dao.InsertTransaction(h.Db, transaction)

		log.Printf("transaction %s applied", transaction.TxId)
	} else {
		log.Printf("transaction %s skipped", transaction.TxId)
	}
}
