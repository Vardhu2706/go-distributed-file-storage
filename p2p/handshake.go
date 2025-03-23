package p2p

// HandshakeFunc defines the signature for a handshake function
// used during peer connection setup. You can customize this to
// perform any kind of peer authentication or negotiation logic.
type HandshakeFunc func(Peer) error

// NOPHandshakeFunc is a no-op handshake function.
// It simply accepts all connections without doing anything.
// Useful for testing or when no handshake validation is required.
func NOPHandshakeFunc(Peer) error {
	return nil
}
