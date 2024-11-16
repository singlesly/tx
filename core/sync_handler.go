package core

import (
	"fmt"
	"github.com/dgraph-io/badger"
	"log"
	"transaction/core/dao"
	"transaction/p2p"
	"transaction/proto"
)

type SyncHandler struct {
	Db   *badger.DB
	Peer p2p.NodePeer
}

func NewSyncHandler(peer p2p.NodePeer, db *badger.DB) *SyncHandler {
	handler := &SyncHandler{Db: db, Peer: peer}

	peer.Subscribe(fmt.Sprintf("node/%s/sync", peer.Node.ID().String()), func(message []byte) {
		var transaction proto.Transaction
		p2p.DeserializeProtoMessage(message, &transaction)

		dao.InsertTransaction(db, &transaction)

		log.Printf("sync transaction %s", transaction.TxId)
	})

	peer.Subscribe(fmt.Sprintf("request-sync"), func(message []byte) {
		go func() {
			db.View(func(txn *badger.Txn) error {
				opts := badger.DefaultIteratorOptions
				opts.Prefix = []byte(dao.TransactionTablePrefix)
				opts.PrefetchSize = 1000
				iterator := txn.NewIterator(opts)
				defer iterator.Close()

				topic, err := peer.PubSub.Join(fmt.Sprintf("node/%s/sync", message))
				if err != nil {
					return err
				}

				for iterator.Rewind(); iterator.Valid(); iterator.Next() {
					tx := dao.ItemToProtoTransaction(iterator.Item())

					peer.Publish(topic, p2p.SerializeProtoMessage(tx))
				}

				return nil
			})
		}()
	})

	return handler
}

func (h *SyncHandler) RequestSync() {
	h.Peer.Publish(h.Peer.Topics["request-sync"], []byte(h.Peer.Node.ID().String()))
}
