package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"strings"
	"time"
	"transaction/p2p"
	"transaction/proto"
)

func main() {

	// Определение флагов для addr, port и списка адресов
	addr := flag.String("addr", "localhost", "Адрес сервера")
	port := flag.Int("port", 8080, "Порт сервера")
	peerList := flag.String("peers", "", "Список пиров в формате address:port, разделённый запятой")

	// Разбор аргументов командной строки
	flag.Parse()

	// Печать значений аргументов
	fmt.Printf("Адрес: %s\n", *addr)
	fmt.Printf("Порт: %d\n", *port)
	fmt.Printf("Список пиров: %s\n", strings.Split(*peerList, ","))

	var peers []string

	if *peerList == "" {
		peers = []string{}
	} else {
		peers = strings.Split(*peerList, ",")
	}

	host := p2p.NewPeer(*addr, *port, peers)

	log.Println("Peer addr", host.Node.Addrs())
	log.Println("Peer id", host.Node.ID())

	topic := host.Subscribe("transactions", func(message []byte) {

		var transaction proto.Transaction
		p2p.DeserializeProtoMessage(message, &transaction)

		log.Println("New message sender: " + transaction.GetSender())
	})

	test := &proto.Transaction{
		TxId:      "123",
		Sender:    "123",
		Recipient: "123",
		Amount:    &proto.Uint256{Value: make([]byte, 3123123)},
		Timestamp: &timestamppb.Timestamp{Seconds: 1628698981},
	}

	<-time.After(10 * time.Second)

	pubData := p2p.SerializeProtoMessage(test)
	var t proto.Transaction
	p2p.DeserializeProtoMessage(pubData, &t)

	go func() {
		for {

			host.Publish(topic, pubData)

			<-time.After(10 * time.Second)
		}
	}()

	select {}
}

// generateBytes генерирует срез байтов указанного размера.
func generateBytes(size int) []byte {
	// Создаем срез нужного размера.
	bytes := make([]byte, size)

	// Заполняем срез случайными байтами.
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}

	return bytes
}
