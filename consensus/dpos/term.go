package dpos

import (
	"encoding/binary"
	"errors"
	"fmt"
	hash2 "github.com/jhdriver/UWORLD/common/hasharry"
	"github.com/jhdriver/UWORLD/consensus"
	"github.com/jhdriver/UWORLD/core/types"
	"github.com/jhdriver/UWORLD/crypto/hash"
	log "github.com/jhdriver/UWORLD/log/log15"
	"github.com/jhdriver/UWORLD/param"
	"math/rand"
	"sort"
)

var errElected = errors.New("has been elected")

// Term calculation for the semester, after each cycle,
// you need to re-elect new super nodes from the candidates
type Term struct {
	dPosStorage IDPosStorage
	chain       consensus.IChain
	time        uint64
}

// Calculate the votes of all candidates.An address can only vote for one
// address, and the real-time balance of the address is used as the number
// of votes it voted for another address.
func (term *Term) countVote() ([]*types.Candidate, error) {
	candidates, err := term.dPosStorage.GetCandidates()
	if err != nil {
		return nil, errors.New("no candidates")
	}
	if len(candidates.Members) < param.MaxWinnerSize {
		return nil, errors.New("too few candidates")
	}
	for index, candidate := range candidates.Members {
		voters := term.dPosStorage.GetCandidateVoters(candidate.Signer)
		for _, voter := range voters {
			candidates.Members[index].Weight += term.chain.GetAddressVote(voter)
		}
	}
	return candidates.Members, nil
}

// Conduct an election
func (term *Term) elect(current uint64, parentHash hash2.Hash, isSort bool) error {
	currentTerm := current / param.TermInterval
	voters, err := term.countVote()
	if err != nil {
		return err
	}
	candidates := types.SortableCandidates{}
	for _, candidate := range voters {
		candidates = append(candidates, candidate)
	}
	if len(candidates) < param.SafeSize {
		return errors.New("too few candidates")
	}

	//
	if isSort {
		sort.Sort(candidates)
	}

	if len(candidates) > param.MaxWinnerSize {
		candidates = candidates[:param.MaxWinnerSize]
	}

	// Use the last block hash of the last cycle as a random number seed
	// to ensure that the election results of each node are consistent
	seed := int64(binary.LittleEndian.Uint32(hash.Hash(parentHash.Bytes()).Bytes())) + int64(currentTerm)
	r := rand.New(rand.NewSource(seed))
	for i := len(candidates) - 1; i > 0; i-- {
		j := int(r.Int31n(int32(i + 1)))
		candidates[i], candidates[j] = candidates[j], candidates[i]
	}

	winners := &types.Winners{Candidates: candidates, ElectParentHash: parentHash}
	term.dPosStorage.SetTermWinners(currentTerm, winners)
	return nil
}

func (term *Term) electCheckTime(chain consensus.IChain, parent, current uint64) error {
	currentTerm := current / param.TermInterval

	winners, _ := term.dPosStorage.GetTermWinners(currentTerm)
	if winners != nil && len(winners.Candidates) != 0 {
		return errElected
	}
	return nil
}

// If a super node has fewer blocks in the previous cycle,
// the super node is kicked out of the candidate and cannot
// participate in subsequent blocks unless it becomes a candidate again
func (term *Term) kickOutValidator(preTerm uint64) error {
	winners, err := term.dPosStorage.GetTermWinners(preTerm)
	if err != nil {
		return fmt.Errorf("failed to get validator: %s", err)
	}
	if len(winners.Candidates) == 0 {
		return errors.New("no winner could be kick out")
	}

	needKickOutWinners := types.SortableCandidates{}
	for _, winner := range winners.Candidates {
		cnt, err := term.dPosStorage.GetTermWinnerMintCnt(preTerm, winner.Signer)
		if err != nil {
			cnt = 0
		}
		if cnt < param.TermInterval/param.BlockInterval/param.MaxWinnerSize/3 {
			needKickOutWinners = append(needKickOutWinners, &types.Candidate{Signer: winner.Signer, PeerId: "", Weight: cnt})
		}
	}
	needKickOutWinnerCnt := len(needKickOutWinners)
	if needKickOutWinnerCnt <= 0 {
		return nil
	}
	sort.Sort(sort.Reverse(needKickOutWinners))
	candidates, err := term.dPosStorage.GetCandidates()
	candidateCount := len(candidates.Members)

	for _, winner := range needKickOutWinners {
		// If the number of candidates is already the smallest, don’t kick out
		if candidateCount <= param.MaxWinnerSize {
			//log.Info("No more candidate can be kic out", "prevEpochID", preTerm, "candidateCount", candidateCount, "needKicOutCount", len(needKickOutWinners)-i)
			return nil
		}
		term.dPosStorage.DeleteCandidate(winner)
		candidateCount--
		log.Info("Kick out candidate", "prevTerm", preTerm, "candidate", winner.Signer.String(), "mintCnt", winner.Weight)
	}
	return nil
}
