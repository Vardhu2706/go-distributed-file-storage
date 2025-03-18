package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPPeer represents the remote node over a TCP established connection.
type TCPPeer struct {

	// conn is the underlying connection of the peer
	conn net.Conn

	// If we dial and retrieve a conn => outbound == True
	// If we accept and retrieve a conn => outbound == False
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn: conn,
		outbound: outbound,
	}
}

type TCPTransport struct {
	listenAddress string
	listener net.Listener
	shakehands HandshakeFunc
	decoder Decoder

	mu sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		shakehands: NOPHandshakeFunc,
		listenAddress: listenAddr,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.listenAddress)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()
	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}

		go t.handleConn(conn)
	}
}

type Temp struct {

}

func (t *TCPTransport) handleConn(conn net.Conn) {
	if err := t.shakehands(conn); err != nil {
		
	}

	// Read Loop
	msg := &Temp{}
	for {
		if err := t.decoder.Decode(conn, msg); err != nil {
			fmt.Printf("TCP error: %s\n", err)
			continue
		}
	}
}