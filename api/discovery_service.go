package api

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"transaction/p2p"
	"transaction/proto"
)

type DiscoveryServiceServer struct {
	proto.UnimplementedDiscoveryServiceServer
	Peer p2p.NodePeer
}

func (s *DiscoveryServiceServer) GetPeers(ctx context.Context, _ *emptypb.Empty) (*proto.GetPeersResponse, error) {
	peers := make([]*proto.Peer, 0)

	for _, item := range s.Peer.Dht.RoutingTable().ListPeers() {
		foundPeer, err := s.Peer.Dht.FindPeer(context.Background(), item)
		if err != nil {
			return nil, err
		}

		addrs := make([]string, 0)
		for _, addr := range foundPeer.Addrs {
			addrs = append(addrs, addr.String())
		}

		peers = append(peers, &proto.Peer{
			Addresses: addrs,
			Id:        foundPeer.ID.String(),
		})

	}

	return &proto.GetPeersResponse{
		Peers: peers,
	}, nil
}
