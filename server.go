package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/Vardhu2706/go-distributed-file-storage/p2p"
)

type FileServerOpts struct {
	StorageRoot 		string
	PathTransformFunc 	PathTransformFunc
	Transport 			p2p.Transport
	BootstrapNodes    []string
}

type FileServer struct {
	FileServerOpts

	peerLock sync.Mutex
	peers map[string]p2p.Peer

	store *Store
	quitch chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root: opts.StorageRoot,
		PathTransformFunc: opts. PathTransformFunc,
	}
	return &FileServer{
		FileServerOpts: opts,
		store:	NewStore(storeOpts),
		quitch: make(chan struct{}),
		peers:	make(map[string]p2p.Peer),
	}
}

func (s *FileServer) broadcast(msg *Message) error {
	peers := []io.Writer{}
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)
	return gob.NewEncoder(mw).Encode(msg)
}

type Message struct {
	Payload any
}

type MessageStoreFile struct {
	Key string
	Size int64
}

func (s *FileServer) StoreData(key string, r io.Reader) error {
	// 1. Store this file to disk.
	// 2. Broadcast this file to all known peers in the network.

	// buf := new(bytes.Buffer)
	// tee := io.TeeReader(r, buf)
	
	// if err := s.store.Write(key, tee); err != nil {
	// 	return err
	// }

	// p := &DataMessage{
	// 	Key: key,
	// 	Data: buf.Bytes(),
	// }

	// fmt.Println(buf.Bytes())
	// return s.broadcast(&Message{
	// 	From: "todo",
	// 	Payload: p,
	// })


	// Testing

	buf := new(bytes.Buffer)

	msg := Message {
		Payload: MessageStoreFile{
			Key: key,
			Size: 15,
		},
	}
	
	if err := gob.NewEncoder(buf).Encode(msg); err != nil {
		return err
	}

	for _, peer := range s.peers {
		if err := peer.Send(buf.Bytes()); err != nil {
			return err
		}
	}


	time.Sleep(time.Second * 2)

	payload := []byte("THIS LARGE FILE")
	for _, peer := range s.peers {
		n, err := io.Copy(peer, bytes.NewReader(payload))
		if err != nil {
			return err
		}

		fmt.Println("Received & Writter bytes to disk: ", n)
	}

	return nil
}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.peers[p.RemoteAddr().String()] = p

	log.Printf("Connected with remote -> %s", p.RemoteAddr())

	return nil
}

func (s *FileServer) loop() {

	defer func ()  {
		log.Println("File server stopped due to user quit action")
		s.Transport.Close()
	}()

	for {
		select {
		case rpc := <-s.Transport.Consume():
			var msg Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				log.Panicln(err)
				return
			}
			
			if err := s.handleMessage(rpc.From, &msg); err != nil {
				log.Panicln(err)
				return
			}
		case <-s.quitch:
			return
		}
	}
}

func (s *FileServer) handleMessage(from string, msg *Message) error {

	switch v := msg.Payload.(type) {
	case MessageStoreFile: 
		return s.handleMessageStoreFile(from, v)
	}

	return nil
}

func (s *FileServer) handleMessageStoreFile(from string, msg MessageStoreFile) error {
	peer, ok := s.peers[from]

	if !ok {
		return fmt.Errorf("peer (%s) not found in the peer list", from)
	}

	if err := s.store.Write(msg.Key, io.LimitReader(peer, msg.Size)); err != nil {
		return err
	}

	peer.(*p2p.TCPPeer).Wg.Done()

	return nil
}

func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.BootstrapNodes {

		if len(addr) == 0 {
			continue
		}

		go func(addr string){
			fmt.Println("Attempting to connect with remote ->", addr)
		
			if err := s.Transport.Dial(addr); err != nil {
				log.Println("Dial Error: ", err)
			}
		}(addr)
	}

	return nil
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	s.bootstrapNetwork()

	s.loop()

	return nil
}


func init() {
	gob.Register(MessageStoreFile{})
}
