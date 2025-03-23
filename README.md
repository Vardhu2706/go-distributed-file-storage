# Go Distributed File Storage

A simplified distributed file storage system written in Go that uses a peer-to-peer (P2P) architecture to share, replicate, and retrieve files across networked nodes.

---

## 💼 Features

- 🚀 Peer-to-peer TCP-based transport layer
- 🧠 Gob-encoded message system + raw stream transfer
- 📁 File storage using content-addressable paths (SHA-1 hashing)
- 📤 Broadcast files to connected peers
- 📥 Request and stream missing files from other nodes
- 🧵 Concurrent processing using goroutines and sync primitives
- 🔎 Manual peer bootstrapping with automatic handshake

---

## 🧱 Architecture

- **Transport Layer**: Handles TCP socket connections, peer handshakes, and message decoding via a custom protocol.
- **RPC System**: Uses `Message` and `RPC` structs to separate high-level and low-level messaging logic.
- **File Server**: Manages connected peers, file storage, message handling, and stream coordination.
- **Store**: Uses SHA-1 hashing to create content-addressable directories and stores files under nested paths.

---

## 🔄 File Flow Example

1. Node A stores a file (`StoreData`)
2. It broadcasts a `MessageStoreFile` to all known peers
3. Other nodes register the file and can fetch it if missing
4. Node B requests the file → Node A responds with a stream
5. File is transferred and written to disk

---

## 🚀 Getting Started

```bash
# Clone this repo
git clone https://github.com/Vardhu2706/go-distributed-file-storage.git
cd go-distributed-file-storage

# Build and run multiple nodes in separate terminals
make run
```

---

## 🧪 Testing

Run transport tests:

```bash
go test ./p2p
```

---

## 🔮 Future Improvements

- File chunking and streaming large files
- Peer reconnection and retry logic
- File deduplication and versioning
- CLI or web interface for uploading & querying files
- Dynamic node discovery (e.g., via DHT or mDNS)