package node

import "fmt"

type PeerNode struct {
	IP          string `json:"ip"`
	Port        uint64 `json:"port"`
	IsBootstrap bool   `json:"is_bootstrap"`
	IsActive    bool   `json:"is_active"`
}

func (pn *PeerNode) tcpAddress() string {
	return fmt.Sprintf("%s:%d", pn.IP, pn.Port)
}

func (pn *PeerNode) apiProtocol() string {
	// Hardcode `apiProtocol` to "http" for now
	return "http"
}

func NewPeerNode(ip string, port uint64, isBootstrap bool, isActive bool) *PeerNode {

	return &PeerNode{IP: ip, Port: port, IsBootstrap: isBootstrap, IsActive: isActive}
}
