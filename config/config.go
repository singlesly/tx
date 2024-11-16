package config

import (
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p/core/peer"
	"strings"
)

type Peer struct {
	Addr   string
	Port   string
	PeerId peer.ID
}

type Config struct {
	Addr    string
	Port    int
	Peers   []Peer
	DataDir string
}

func NewConfigFromCommandArgs() *Config {
	// Определение флагов для addr, port и списка адресов
	addr := flag.String("addr", "0.0.0.0", "Адрес сервера")
	port := flag.Int("port", 0, "Порт сервера")
	peerList := flag.String("peers", "", "Список пиров в формате address:port, разделённый запятой")
	dataDir := flag.String("data-dir", "./var/data", "Путь к базе данных")
	// Разбор аргументов командной строки
	flag.Parse()

	// Печать значений аргументов
	fmt.Printf("Адрес: %s\n", *addr)
	fmt.Printf("Порт: %d\n", *port)
	fmt.Printf("Список пиров: %s\n", strings.Split(*peerList, ","))
	fmt.Printf("Путь к базе данных: %s\n", *dataDir)

	var peers []Peer

	if *peerList == "" {
		peers = []Peer{}
	} else {
		peersStrs := strings.Split(*peerList, ",")
		for _, peerStr := range peersStrs {
			splittedPeer := strings.Split(peerStr, ":")
			pAddr := splittedPeer[0]
			pPort := splittedPeer[1]

			peers = append(peers, Peer{
				Addr:   pAddr,
				Port:   pPort,
				PeerId: "faked-id",
			})
		}
	}

	return &Config{
		Addr:    *addr,
		Port:    *port,
		Peers:   peers,
		DataDir: *dataDir,
	}
}
