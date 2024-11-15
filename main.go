package main

import (
	"log"
	config2 "transaction/config"
	"transaction/p2p"
)

func main() {
	config := config2.NewConfigFromCommandArgs()

	host := p2p.NewPeer(config.Addr, config.Port, config.Peers)

	log.Printf("Peer addr=%s, id=%s", host.Node.Addrs(), host.Node.ID())

	select {}
}
