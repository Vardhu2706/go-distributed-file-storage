package p2p

// Constants that mark the type of incoming RPC message.
// These are used by the decoder to distinguish between
// regular gob-encoded messages and raw stream data.
const (
	IncomingMessage = 0x1 // Identifies a standard message (e.g., gob-encoded)
	IncomingStream  = 0x2 // Identifies a stream message (e.g., raw file transfer)
)

// RPC holds any arbitrary data that is being sent over
// each transport between two nodes in the network.
type RPC struct {
	// From represents the address of the peer that sent the message
	From string

	// Payload holds the raw byte content of the message
	Payload []byte

	// Stream indicates whether the message is a stream (true) or regular (false)
	Stream bool
}
