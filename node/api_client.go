package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tejasa97/go-block/database"
	"net/http"
)

const (
	endpointNodeStatus      = "/node/status"
	endpointNodeBlocksQuery = "/node/blocks"
	endpointNodeAdd         = "/node/addPeer"
	endpointBalancesList    = "/balances/list"
)

type ApiClient interface {
	getPeerStatus(peer PeerNode) (StateRes, error)
	addPeer(node Node, peer PeerNode) (StateRes, error)
	queryBlocks(peerNode PeerNode, blockHash database.Hash) ([]database.Block, error)
}

type apiClient struct {
}

func NewApiClient() ApiClient {
	return &apiClient{}
}

func (n *apiClient) getPeerStatus(peerNode PeerNode) (StateRes, error) {

	url := fmt.Sprintf("%s://%s%s", peerNode.apiProtocol(), peerNode.tcpAddress(), endpointNodeStatus)
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

func (n *apiClient) addPeer(node Node, peerNode PeerNode) (StateRes, error) {

	url := fmt.Sprintf("%s://%s%s", peerNode.apiProtocol(), peerNode.tcpAddress(), endpointNodeAdd)
	data, err := json.Marshal(node.info)

	res, err := http.Post(url, "application/json", bytes.NewBuffer(data))
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

func (n *apiClient) queryBlocks(peerNode PeerNode, blockHash database.Hash) ([]database.Block, error) {

	url := fmt.Sprintf("%s://%s%s?fromHash=%s", peerNode.apiProtocol(), peerNode.tcpAddress(), endpointNodeBlocksQuery, blockHash.Hex())
	res, err := http.Get(url)
	if err != nil {
		return []database.Block{}, err
	}

	blocksRes := BlocksRes{}
	err = readRes(res, &blocksRes)
	if err != nil {
		return []database.Block{}, err
	}

	return blocksRes.Blocks, nil
}
