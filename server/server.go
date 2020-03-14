package server

import (
	"fmt"
	"log"
	"net"

	"github.com/mendersoftware/go-lib-micro/config"

	dconfig "github.com/canyanio/rating-agent-hep/config"
)

// StartServer starts the UDP server which receives the HEP packats
func StartServer() error {
	listen := config.Config.GetString(dconfig.SettingListen)
	log.Printf("Listening on %v", listen)

	pc, err := net.ListenPacket("udp", listen)
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	for {
		buf := make([]byte, 65536)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			continue
		}
		go serve(pc, addr, buf[:n])
	}
}

func serve(pc net.PacketConn, addr net.Addr, buf []byte) {
	fmt.Println(string(buf))
	pc.WriteTo(buf, addr)
}
