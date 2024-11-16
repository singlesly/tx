package core

import (
	"github.com/dgraph-io/badger"
	"transaction/p2p"
)

type Core struct{}

func NewCore(peer p2p.NodePeer, db *badger.DB) {
	NewTransactionHandler(peer, db)
	syncHandler := NewSyncHandler(peer, db)
	syncHandler.RequestSync()
}
