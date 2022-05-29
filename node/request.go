package node

type TxAddReq struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value uint   `json:"value"`
	Data  string `json:"data"`
}

type AddNodeReq struct {
	PeerNode
}

func (addNodeReq AddNodeReq) getPeerNode() PeerNode {
	return PeerNode{
		IP:          addNodeReq.IP,
		Port:        addNodeReq.Port,
		IsBootstrap: addNodeReq.IsBootstrap,
		IsActive:    addNodeReq.IsActive,
	}
}
