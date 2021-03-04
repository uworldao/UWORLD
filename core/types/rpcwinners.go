package types

import (
	"github.com/uworldao/UWORLD/common/hasharry"
)

type RpcWinners struct {
	Candidates      []*RpcCandidate `json:"winners"`
	ElectParentHash string          `json:"electblockhash"`
}

func TranslateWinnersToRpcWinners(winners *Winners, mntCount map[hasharry.Address]uint64) *RpcWinners {
	rpcWinners := &RpcWinners{
		Candidates:      make([]*RpcCandidate, 0),
		ElectParentHash: "",
	}
	for _, winner := range winners.Candidates {
		rpcCandidate := &RpcCandidate{
			Signer:   winner.Signer.String(),
			PeerId:   winner.PeerId,
			Weight:   winner.Weight,
			MntCount: mntCount[winner.Signer],
		}
		rpcWinners.Candidates = append(rpcWinners.Candidates, rpcCandidate)
	}
	rpcWinners.ElectParentHash = winners.ElectParentHash.String()
	return rpcWinners
}
