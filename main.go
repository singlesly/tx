package main

import (
	"github.com/dgraph-io/badger"
	"log"
	"transaction/api"
	config2 "transaction/config"
	"transaction/core"
	"transaction/p2p"
)

func main() {
	config := config2.NewConfigFromCommandArgs()

	host := p2p.NewPeer(config.Addr, config.Port, config.Peers)

	log.Printf("Peer addr = %s, id = %s", host.Node.Addrs(), host.Node.ID())

	db, err := badger.Open(badger.DefaultOptions(config.DataDir))

	if err != nil {
		log.Printf("cannot connect db")
		panic(err)
	}

	core.NewCore(host, db)
	api.NewApi(host, db, config)

	select {}
}
