package p2p

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	peer2 "github.com/libp2p/go-libp2p/core/peer"
	"github.com/mr-tron/base58"
	"github.com/multiformats/go-multiaddr"
	"google.golang.org/protobuf/proto"
	"log"
	"strings"
)

type NodePeer struct {
	Node   host.Host
	PubSub *pubsub.PubSub
	topics map[string]*pubsub.Topic
}

func NewPeer(addr string, port int, peerList []string) NodePeer {
	node, err := libp2p.New(libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/%s/tcp/%d", addr, port)))

	if err != nil {
		log.Fatalf(err.Error())
		panic(err)
	}

	for _, peerToConnect := range peerList {
		ConnectPeer(&node, peerToConnect, "")
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

	topics := make(map[string]*pubsub.Topic)

	return NodePeer{Node: node, PubSub: ps, topics: topics}
}

func ConnectPeer(node *host.Host, addr string, pid peer2.ID) {
	splittedPeer := strings.Split(addr, ":")
	pAddr := splittedPeer[0]
	pPort := splittedPeer[1]
	pfAddress := fmt.Sprintf("/ip4/%s/tcp/%s", pAddr, pPort)
	pMultiAddr, pErr := multiaddr.NewMultiaddr(pfAddress)
	if pErr != nil {
		panic(pErr)
	}

	if pid == "" {
		pid = "faked"
	}

	pInfo := peer2.AddrInfo{
		ID:    pid,
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

			ConnectPeer(node, addr, actualPidBytes)

			return
		}
		panic(cErr)
	} else {
		log.Printf("connected to peer %s", addr)
	}
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
	publishError := topic.Publish(context.Background(), compressData(serializedMessage))

	if publishError != nil {
		log.Printf("Ошибка при публикации сообщения: %v", publishError)
	}

	log.Println("Publish message")
}

func (np *NodePeer) Subscribe(topicName string, callback func(message []byte)) *pubsub.Topic {
	topic, joinErr := np.PubSub.Join(topicName)

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

			if msg.GetFrom().String() == np.Node.ID().String() {
				continue
			}

			log.Printf("Incoming message from node: %s", msg.GetFrom().String())

			callback(decompressData(msg.GetData()))
		}
	}()

	return topic
}
