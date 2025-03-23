package p2p

import (
	"encoding/gob"
	"io"
)

// Decoder defines an interface for decoding RPC messages from an io.Reader.
// Different decoding strategies (e.g., gob vs custom binary protocol) can implement this.
type Decoder interface {
	Decode(io.Reader, *RPC) error
}

// GOBDecoder uses Go's built-in gob decoder to deserialize data into an RPC struct.
type GOBDecoder struct{}

// Decode implements the Decoder interface using gob.
func (dec GOBDecoder) Decode(r io.Reader, msg *RPC) error {
	return gob.NewDecoder(r).Decode(msg)
}

// DefaultDecoder is a custom decoder that handles two cases:
// - If the message is a stream indicator, mark RPC as a stream.
// - Otherwise, read raw bytes into the RPC payload.
type DefaultDecoder struct{}

// Decode reads 1 byte to peek the message type, then either sets Stream flag or loads a basic payload.
func (dec DefaultDecoder) Decode(r io.Reader, msg *RPC) error {
	peekBuf := make([]byte, 1)
	if _, err := r.Read(peekBuf); err != nil {
		return nil // silently return, assuming stream not ready
	}

	stream := peekBuf[0] == IncomingStream
	if stream {
		msg.Stream = true
		return nil
	}

	// Otherwise, treat the rest of the message as a raw byte payload.
	buf := make([]byte, 1028)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}

	msg.Payload = buf[:n]
	return nil
}
