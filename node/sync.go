package node

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func (n *Node) sync(ctx context.Context) error {
	ticker := time.NewTicker(2 * time.Second)

	for {
		select {
		case <-ticker.C:
			fmt.Println("Searching for new Peers and blocks..")
			n.fetchNewBlocksAndPeers()

		case <-ctx.Done():
			ticker.Stop()
		}
	}
	return nil
}

func (n *Node) fetchNewBlocksAndPeers() error {

	for _, peer := range n.knownPeers {
		// Same node
		if n.info.IP == peer.IP && n.info.Port == peer.Port {
			continue
		}

		nodeState, err := getPeerStatus(peer)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}

		localBlockNumber := n.state.GetLatestBlockHeader().Number
		if localBlockNumber < nodeState.BlockNumber {
			newBlocksCount := nodeState.BlockNumber - localBlockNumber
			fmt.Printf("Found %d new blocks from Peer %s\n", newBlocksCount, peer.IP)
		}

		for _, peerNode := range nodeState.KnownPeers {
			newPeer, isKnownPeer := n.knownPeers[peerNode.tcpAddress()]
			if !isKnownPeer {
				fmt.Printf("Found new peer %s \n", peer.tcpAddress())
				n.knownPeers[peerNode.tcpAddress()] = newPeer
			}
		}
	}

	return nil
}

func getPeerStatus(peer PeerNode) (StateRes, error) {
	url := fmt.Sprintf("%s://%s%s", peer.apiProtocol(), peer.tcpAddress(), endpointNodeStatus)
	res, err := http.Get(url)
	if err != nil {
		return StateRes{}, err
	}

	nodeStateRes := StateRes{}
	err = readRes(res, &nodeStateRes)
	if err != nil {
		return StateRes{}, err
	}

	return nodeStateRes, nil
}
