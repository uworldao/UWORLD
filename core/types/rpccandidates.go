package types

type RpcCandidate struct {
	Signer   string `json:"address"`
	PeerId   string `json:"peerid"`
	Weight   uint64 `json:"votes"`
	MntCount uint64 `json:"mntcount"`
}

type RpcCandidates struct {
	Candidates []*RpcCandidate `json:"candidates"`
}

func TranslateCandidatesToRpcCandidates(candidates []*Candidate) *RpcCandidates {
	rpcCandidates := &RpcCandidates{Candidates: make([]*RpcCandidate, 0)}
	for _, candidate := range candidates {
		rpcCandidate := &RpcCandidate{
			Signer: candidate.Signer.String(),
			PeerId: candidate.PeerId,
			Weight: candidate.Weight,
		}
		rpcCandidates.Candidates = append(rpcCandidates.Candidates, rpcCandidate)
	}
	return rpcCandidates
}
