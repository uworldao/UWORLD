package types

// Node information
type NodeInfo struct {
	// Node version
	Version string `json:"version"`
	// Node network
	Net string `json:"net"`
	// Node p2p id
	P2pId string `json:"p2pid"`
	// Node p2p address
	P2pAddr string `json:"p2pAddr"`
	// Linked node
	Connections uint32 `json:"connections"`
	// Current block height
	Height uint64 `json:"height"`
	// Current effective block height
	Confirmed uint64 `json:"confirmed"`
}
