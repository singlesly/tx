package main

import (
	"github.com/dgraph-io/badger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"transaction/api"
	config2 "transaction/config"
	"transaction/core"
	"transaction/p2p"
	"transaction/proto"
)

func main() {
	config := config2.NewConfigFromCommandArgs()

	host := p2p.NewPeer(config.Addr, config.Port, config.Peers)

	log.Printf("Peer addr = %s, id = %s", host.Node.Addrs(), host.Node.ID())

	db, err := badger.Open(badger.DefaultOptions(config.DataDir))

	core.NewTransactionHandler(host, db)

	if err != nil {
		log.Printf("cannot connect db")
	}

	go func() {
		lis, _ := net.Listen("tcp", ":50501")
		s := grpc.NewServer()
		proto.RegisterTransactionServiceServer(s, &api.TransactionServiceServer{
			Db:   db,
			Peer: host,
		})
		reflection.Register(s)
		log.Printf("listen 0.0.0.0:50501")

		sErr := s.Serve(lis)
		if sErr != nil {
			log.Fatalf(sErr.Error())
			return
		}
	}()

	select {}
}
