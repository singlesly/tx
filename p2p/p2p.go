package p2p

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	peer2 "github.com/libp2p/go-libp2p/core/peer"
	"github.com/mr-tron/base58"
	"github.com/multiformats/go-multiaddr"
	"google.golang.org/protobuf/proto"
	"log"
	"strings"
	"transaction/config"
)

type NodePeer struct {
	Node   host.Host
	Dht    *dht.IpfsDHT
	PubSub *pubsub.PubSub
	Topics map[string]*pubsub.Topic
}

func NewPeer(addr string, port int, peers []config.Peer) NodePeer {
	node, err := libp2p.New(libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/%s/tcp/%d", addr, port)))

	if err != nil {
		log.Fatalf(err.Error())
		panic(err)
	}

	connectedPeers := make([]peer2.AddrInfo, 0)

	for _, peerToConnect := range peers {
		connectedPeers = ConnectPeer(&node, peerToConnect)
	}

	node.Network().Notify(&network.NotifyBundle{
		ConnectedF: func(n network.Network, conn network.Conn) {
			peerId := conn.RemoteMultiaddr().String()
			log.Printf("new peer found %s", peerId)
		},
	})

	ps, psErr := pubsub.NewGossipSub(context.Background(), node)

	if psErr != nil {
		panic(psErr)
	}

	ipfsDht := initDht(node, connectedPeers)

	return NodePeer{Node: node, Dht: ipfsDht, PubSub: ps, Topics: make(map[string]*pubsub.Topic)}
}

func initDht(node host.Host, peers []peer2.AddrInfo) *dht.IpfsDHT {

	ipfsDHT, err := dht.New(context.Background(), node, dht.Mode(dht.ModeServer))
	if err != nil {
		panic("cannot init dht")
	}

	bErr := ipfsDHT.Bootstrap(context.Background())
	if bErr != nil {
		panic("cannot bootstrap dht")
	}

	return ipfsDHT
}

func ConnectPeer(node *host.Host, peer config.Peer) []peer2.AddrInfo {

	connectedPeers := make([]peer2.AddrInfo, 0)

	pfAddress := fmt.Sprintf("/ip4/%s/tcp/%s", peer.Addr, peer.Port)
	pMultiAddr, pErr := multiaddr.NewMultiaddr(pfAddress)
	if pErr != nil {
		panic(pErr)
	}

	pInfo := peer2.AddrInfo{
		ID:    peer.PeerId,
		Addrs: []multiaddr.Multiaddr{pMultiAddr},
	}

	log.Println(pInfo)

	cErr := (*node).Connect(context.Background(), pInfo)

	if cErr != nil {
		if strings.Contains(cErr.Error(), "peer id mismatch:") {
			log.Println("found mismatch")
			errorParts := strings.Split(cErr.Error(), " ")
			actualPid := errorParts[len(errorParts)-1]

			actualPidDecoded, _ := base58.Decode(actualPid)
			actualPidBytes, _ := peer2.IDFromBytes(actualPidDecoded)

			pInfo := peer2.AddrInfo{
				ID:    actualPidBytes,
				Addrs: []multiaddr.Multiaddr{pMultiAddr},
			}

			log.Println(pInfo)

			ConnectPeer(node, config.Peer{
				Addr:   peer.Addr,
				Port:   peer.Port,
				PeerId: actualPidBytes,
			})

			connectedPeers = append(connectedPeers, pInfo)

			return connectedPeers
		}
		panic(cErr)
	} else {
		log.Printf("connected to peer %s", peer.Addr)
	}

	return connectedPeers
}

func SerializeProtoMessage[T proto.Message](message T) []byte {
	msg, err := proto.Marshal(message)

	if err != nil {
		panic(err)
	}

	return msg
}

func DeserializeProtoMessage[T proto.Message](message []byte, result T) {
	err := proto.Unmarshal(message, result)
	if err != nil {
		panic(err)
	}
}

func (np *NodePeer) Publish(topic *pubsub.Topic, serializedMessage []byte) {
	msg := compressData(serializedMessage)
	msgSizeError := CheckMessageSize(msg)

	if msgSizeError != nil {
		log.Println(msgSizeError.Error())
		return
	}

	publishError := topic.Publish(context.Background(), compressData(serializedMessage))

	if publishError != nil {
		log.Printf("error on publish message: %v", publishError)
	}

	log.Printf("[%s] publish with size %d", topic.String(), len(msg))
}

func (np *NodePeer) Subscribe(topicName string, callback func(message []byte)) *pubsub.Topic {
	topic, joinErr := np.PubSub.Join(topicName)
	(np.Topics)[topicName] = topic

	log.Printf("subscribed " + topicName)

	x, _ := (np.Topics)[topicName]
	log.Printf(x.String())

	if joinErr != nil {
		panic(joinErr)
	}

	sub, subErr := topic.Subscribe()

	if subErr != nil {
		panic(subErr)
	}

	go func() {
		for {

			msg, err := sub.Next(context.Background())

			if err != nil {
				log.Printf("Ошибка при получении сообщения: %v", err)
				continue
			}

			//if msg.GetFrom().String() == np.Node.ID().String() {
			//	continue
			//}

			log.Printf("Incoming message from node: %s", msg.GetFrom().String())

			callback(decompressData(msg.GetData()))
		}
	}()

	return topic
}
