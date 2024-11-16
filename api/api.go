package api

import (
	"fmt"
	"github.com/dgraph-io/badger"
	"google.golang.org/grpc"
	"log"
	"net"
	"transaction/config"
	"transaction/p2p"
	"transaction/proto"
)

type Api struct{}

func NewApi(peer p2p.NodePeer, db *badger.DB, config *config.Config) {
	lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", config.GrpcPort))
	s := grpc.NewServer()
	registerServices(peer, db, s)
	log.Printf("listen %s", lis.Addr().String())

	go func() {
		sErr := s.Serve(lis)
		if sErr != nil {
			log.Fatalf(sErr.Error())
			return
		}
	}()
}

func registerServices(peer p2p.NodePeer, db *badger.DB, s *grpc.Server) {
	proto.RegisterTransactionServiceServer(s, &TransactionServiceServer{
		Db:   db,
		Peer: peer,
	})
}
