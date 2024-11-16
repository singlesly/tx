package api

import (
	"context"
	"github.com/dgraph-io/badger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"transaction/core/dao"
	"transaction/p2p"
	"transaction/proto"
)

type TransactionServiceServer struct {
	proto.UnimplementedTransactionServiceServer
	Db   *badger.DB
	Peer p2p.NodePeer
}

func (s *TransactionServiceServer) CreateTransaction(ctx context.Context, req *proto.CreateTransactionRequest) (*emptypb.Empty, error) {
	topic, exists := s.Peer.Topics["transactions"]

	if exists == false {
		return nil, status.Errorf(codes.Internal, "Topic for published not found")
	}

	s.Peer.Publish(topic, p2p.SerializeProtoMessage(req.Transaction))

	return &emptypb.Empty{}, nil
}

func (s *TransactionServiceServer) GetTransaction(ctx context.Context, req *proto.GetTransactionRequest) (*proto.GetTransactionResponse, error) {

	transaction, err := dao.GetTransaction(s.Db, req.TxId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Tx not found")
	}

	return &proto.GetTransactionResponse{
		Transaction: transaction,
	}, nil
}
