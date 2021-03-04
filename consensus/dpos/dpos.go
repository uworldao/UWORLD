package dpos

import (
	"errors"
	"fmt"
	"github.com/uworldao/UWORLD/common/hasharry"
	"github.com/uworldao/UWORLD/consensus"
	"github.com/uworldao/UWORLD/core/types"
	"github.com/uworldao/UWORLD/crypto/hash"
	"github.com/uworldao/UWORLD/database/dposdb"
	"github.com/uworldao/UWORLD/param"
	"sort"
	"time"
)

const (

	// DPos storage file name
	dposStorage = "dpos"
)

type DPos struct {
	dposStorage          IDPosStorage
	signer               hasharry.Address
	sign                 consensus.ISign
	confirmedBlockHeader *types.Header
}

func NewDPos(DataDir string, signer hasharry.Address, sign consensus.ISign) (*DPos, error) {
	dposStorage := dposdb.NewDPosStorage(DataDir + "/" + dposStorage)
	if err := dposStorage.Open(); err != nil {
		return nil, err
	}

	dpos := &DPos{
		dposStorage:          dposStorage,
		signer:               signer,
		sign:                 sign,
		confirmedBlockHeader: nil,
	}
	return dpos, nil
}

func (dpos *DPos) Close() error {
	return dpos.dposStorage.Close()
}

func (dpos *DPos) Init(chain consensus.IChain) error {
	gensis, err := chain.GetHeaderByHeight(0)
	if err != nil {
		return err
	}
	confirmedHash, err := dpos.dposStorage.GetConfirmedBlockHash()
	if err != nil {
		dpos.confirmedBlockHeader = gensis
	} else {
		if dpos.confirmedBlockHeader, err = chain.GetHeaderByHash(confirmedHash); err != nil {
			return err
		}
	}
	return nil //dpos.elect(gensis.Time, gensis.Hash, chain, false)
}

func (dpos *DPos) GetGenesisBlock() *types.Block {
	block := &types.Block{
		Header: &types.Header{
			Hash:          hasharry.Hash{},
			ParentHash:    hasharry.Hash{},
			TxRoot:        hasharry.Hash{},
			StateRoot:     hasharry.Hash{},
			ContractRoot:  hasharry.Hash{},
			ConsensusRoot: hasharry.Hash{},
			Height:        0,
			Time:          1569398062,
			Term:          0,
			SignScript:    &types.SignScript{},
			Signer:        hasharry.Address{},
		},
		Body: &types.Body{Transactions: types.Transactions{}},
	}
	for _, info := range initialCandidates {
		var peerId types.PeerId
		copy(peerId[:], info.PeerId)
		tx := &types.Transaction{
			TxHead: &types.TransactionHead{
				TxHash:     hasharry.Hash{},
				TxType:     types.LoginCandidate,
				From:       hasharry.StringToAddress(info.Address),
				Nonce:      0,
				Fees:       0,
				Time:       1569398062,
				SignScript: &types.SignScript{},
			},
			TxBody: &types.LoginTransactionBody{
				PeerId: peerId,
			},
		}
		tx.SetHash()
		block.Transactions = append(block.Transactions, tx)
	}
	var sumCoins uint64
	for _, info := range param.MappingCoin {
		tx := &types.Transaction{
			TxHead: &types.TransactionHead{
				TxHash:     hasharry.Hash{},
				TxType:     types.NormalTransaction,
				From:       hasharry.StringToAddress(info.Address),
				Nonce:      0,
				Fees:       0,
				Time:       1569398062,
				Note:       info.Note,
				SignScript: &types.SignScript{},
			},
			TxBody: &types.NormalTransactionBody{
				Contract: param.Token,
				To:       hasharry.StringToAddress(info.Address),
				Amount:   info.Amount,
			},
		}
		sumCoins += info.Amount
		tx.SetHash()
		block.Transactions = append(block.Transactions, tx)
	}

	block.TxRoot = block.Transactions.Hash()
	block.Hash = hash.Hash(block.ToBytes())
	return block
}

func (dpos *DPos) GetTermInterval() uint64 {
	return param.TermInterval
}

func (dpos *DPos) Sign(block *types.Block) error {
	var err error
	if block.Height == 0 {
		return errors.New("unknown block")
	}
	block.SignScript, err = dpos.sign.SignHash(block.Hash)
	if err != nil {
		return err
	}
	return nil
}

// Check whether the block header signature address is a block at this time
func (dpos *DPos) CheckWinner(chain consensus.IChain, header *types.Header) error {
	parentHead, err := chain.GetHeaderByHash(header.ParentHash)
	if err != nil {
		return err
	}

	if err := dpos.electCheckTerm(chain, parentHead.Time, header.Time); err != errElected {
		// If elections have not yet taken place in this cycle, conduct elections first
		if err := dpos.elect(header.Time, parentHead.Hash, chain, true); err != nil {
			return err
		}
	}

	// Check if the time of block production is correct
	if err := dpos.checkTime(parentHead, header); err != nil {
		return err
	}

	// Find the block address at that time
	winner, err := dpos.lookupWinners(header.Time)
	if err != nil {
		return err
	}
	if !winner.IsEqual(dpos.signer) {
		return errors.New("it's not the miner's turn")
	}
	return nil
}

// Verify block header time and signature
func (dpos *DPos) VerifyHeader(header, parent *types.Header) error {
	// If the block time is in the future, it will fail
	if header.Time > uint64(time.Now().Unix()) {
		return errors.New("block in the future")
	}
	// Verify whether it is the time point of block generation
	if err := dpos.checkTime(parent, header); err != nil {
		return errors.New("time check failed")
	}
	if header.SignScript == nil {
		return errors.New("no signature")
	}
	if parent.Time+param.BlockInterval > header.Time {
		return errors.New("invalid timestamp")
	}
	return nil
}

func (dpos *DPos) VerifySeal(chain consensus.IChain, header *types.Header, parent *types.Header) error {
	// Verifying the genesis block is not supported
	if header.Height == 0 {
		return errors.New("unknown block")
	}
	if dpos.confirmedBlockHeader != nil {
		if header.Height <= dpos.confirmedBlockHeader.Height {
			return errors.New("height error")
		}
	}
	parent, err := chain.GetHeaderByHash(header.ParentHash)
	if err != nil {
		return errors.New("unknown parent hash")
	}
	preTermLastHeader, err := dpos.getPreTermLastBlockHash(header, chain)
	if err != nil {
		return err
	}
	// Verify the block node
	if err := dpos.VerifyCreator(header, preTermLastHeader, chain); err != nil {
		return err
	}
	// Update the height of the confirmed block
	return dpos.updateConfirmedBlockHeader(chain)
}

// If the current number of candidates is less than or equal to the
// number of super nodes, it is not allowed to withdraw candidates.
func (dpos *DPos) VerifyTx(tx types.ITransaction) error {
	switch tx.GetTxType() {
	/*case types.LogoutCandidate:
	cans, _ := dpos.dposStorage.GetCandidates()
	if cans.Len() <= maxWinnerSize {
		return fmt.Errorf("candidate nodes are already in the minimum number. Cannot cancel the candidate status now, please wait")
	}*/
	}
	return nil
}

// Verify that the address of the block generated at this time is correct,
// and verify the signature.
func (dpos *DPos) VerifyCreator(header *types.Header, parent *types.Header, chain consensus.IChain) error {
	signer, err := dpos.lookupWinnerNoExistToCreate(header.Time, parent, chain)
	if err != nil {
		return err
	}
	if err := dpos.verifyBlockSigner(signer, header); err != nil {
		return err
	}
	return nil
}

func (dpos *DPos) GetWinnersPeerID(time uint64) ([]string, error) {
	term := time / param.TermInterval
	winners, err := dpos.dposStorage.GetTermWinners(term)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, winner := range winners.Candidates {
		ids = append(ids, winner.PeerId)
	}
	return ids, nil
}

func (dpos *DPos) GetConfirmedBlockHeader(chain consensus.IChain) *types.Header {
	if dpos.confirmedBlockHeader == nil {
		header, err := dpos.loadConfirmedBlockHeader(chain)
		if err != nil {
			header, _ = chain.GetHeaderByHeight(0)
			dpos.confirmedBlockHeader = header
			return header
		} else {
			return header
		}
	}
	return dpos.confirmedBlockHeader
}

func (dpos *DPos) GetCandidates(chain consensus.IChain) []*types.Candidate {
	candidates, err := dpos.dposStorage.GetCandidates()
	if err != nil {
		return nil
	}
	for index, candidate := range candidates.Members {
		voters := dpos.dposStorage.GetCandidateVoters(candidate.Signer)
		for _, voter := range voters {
			candidates.Members[index].Weight += chain.GetAddressVote(voter)
		}
	}
	sortedCandidates := types.SortableCandidates{}
	for _, candidate := range candidates.Members {
		sortedCandidates = append(sortedCandidates, candidate)
	}
	sort.Sort(sortedCandidates)
	return sortedCandidates
}

func (dpos *DPos) GetTermWinners(term uint64) *types.Winners {
	winners, _ := dpos.dposStorage.GetTermWinners(term)
	return winners
}

func (dpos *DPos) GetTermWinnersMntCount(term uint64, address hasharry.Address) (uint64, error) {
	return dpos.dposStorage.GetTermWinnerMintCnt(term, address)
}

func (dpos *DPos) SetConfirmedHeader(header *types.Header) {
	dpos.confirmedBlockHeader = header
	dpos.dposStorage.SetConfirmedBlockHash(header.Hash)
}

func (dpos *DPos) InitTrie(consensusRoot hasharry.Hash) error {
	return dpos.dposStorage.InitTrie(consensusRoot)
}

func (dpos *DPos) Commit() (hasharry.Hash, error) {
	return dpos.dposStorage.Commit()
}

func (dpos *DPos) RootHash() hasharry.Hash {
	return dpos.dposStorage.RootHash()
}

func (dpos *DPos) electCheckTerm(chain consensus.IChain, currentTime, now uint64) error {
	term := &Term{dPosStorage: dpos.dposStorage}
	return term.electCheckTime(chain, currentTime, now)
}

func (dpos *DPos) elect(now uint64, preTermLastBlockHash hasharry.Hash, chain consensus.IChain, isSort bool) error {
	term := &Term{dpos.dposStorage, chain, now}
	err := term.elect(now, preTermLastBlockHash, isSort)
	if err != nil {
		return errors.New(fmt.Sprintf("got error when elect next term, err: %s", err))
	}
	return nil
}

func (dpos *DPos) setWinners(time uint64, parent *types.Header, reader consensus.IChain) ([]*types.Candidate, error) {
	term := time / param.TermInterval
	winners, err := dpos.dposStorage.GetTermWinners(term)

	// If the election result of the current cycle does not
	// exist, the current cycle of elections is conducted
	if err != nil || winners == nil || !parent.Hash.IsEqual(winners.ElectParentHash) {
		dpos.electCheckTerm(reader, parent.Time, time)
		if err := dpos.elect(time, parent.Hash, reader, true); err != nil {
			return nil, err
		}
		if winners, err = dpos.dposStorage.GetTermWinners(term); err != nil {
			return nil, err
		}
	}
	return winners.Candidates, nil
}

// updateConfirmedBlockHeader Update the final confirmation block
func (dpos *DPos) updateConfirmedBlockHeader(chain consensus.IChain) error {
	if dpos.confirmedBlockHeader == nil {
		header, err := dpos.loadConfirmedBlockHeader(chain)
		if err != nil {
			header, err = chain.GetHeaderByHeight(0)
			if err != nil {
				return err
			}
		}
		dpos.confirmedBlockHeader = header
	}
	curHeader, err := chain.CurrentHeader()
	if err != nil {
		return err
	}
	// If there are already more than two-thirds of different nodes generating blocks,
	// it means that the blocks before these blocks have been confirmed

	term := uint64(0)
	winnerMap := make(map[string]int)
	for !dpos.confirmedBlockHeader.Hash.IsEqual(curHeader.Hash) &&
		dpos.confirmedBlockHeader.Height < curHeader.Height {
		curTerm := curHeader.Time / param.TermInterval
		if curTerm != term {
			term = curTerm
			winnerMap = make(map[string]int)
		}
		// fast return
		// if block number difference less consensusSize-witnessNum
		// there is no need to check block is confirmed

		count := winnerMap[curHeader.Signer.String()]
		winnerMap[curHeader.Signer.String()] = count + 1

		if len(winnerMap) >= param.ConsensusSize /*dpos.checkWinnerMapCount(winnerMap, 1)*/ {
			dpos.dposStorage.SetConfirmedBlockHash(curHeader.Hash)
			dpos.confirmedBlockHeader = curHeader
			chain.UpdateConfirmedHeight(curHeader.Height)
			//log.Info("DPos set confirmed block header", "currentHeader", curHeader.Height)
			return nil
		}
		curHeader, err = chain.GetHeaderByHash(curHeader.ParentHash)
		if err != nil {
			return errors.New("nil block header returned")
		}
	}
	return nil
}

// If the number of outgoing block nodes is greater than the
// minimum confirmation number, the block is confirmed as valid
func (dpos *DPos) checkWinnerMapCount(winnerMap map[string]int, maxCount int) bool {
	if len(winnerMap) < param.ConsensusSize {
		return false
	}
	winnerCount := 0
	for _, count := range winnerMap {
		if count >= maxCount {
			winnerCount++
		}
		if winnerCount >= param.ConsensusSize {
			return true
		}
	}
	return false
}

func (dpos *DPos) updateGenesisDPosStorage(chain consensus.IChain) error {
	block, err := chain.GetBlockByHeight(0)
	if err != nil {
		return err
	}
	dpos.UpdateConsensus(block)
	return nil
}

// Update consensus candidates and voting information
func (dpos *DPos) UpdateConsensus(block *types.Block) {
	for _, tx := range block.Transactions {
		switch tx.GetTxType() {
		case types.LoginCandidate:
			candidate := &types.Candidate{
				Signer: tx.From(),
				PeerId: string(tx.GetTxBody().GetPeerId()),
				Weight: 0,
			}
			// When becoming a candidate, also vote for yourself
			dpos.dposStorage.SetCandidate(candidate)
			dpos.dposStorage.SetVoter(tx.From(), tx.From())
			/*case types.LogoutCandidate:
				candidate := &types.Candidate{
					Signer: tx.From(),
					PeerId: "",
					Weight: 0,
				}
				dpos.dposStorage.DeleteCandidate(candidate)
			case types.VoteToCandidate:
				dpos.dposStorage.SetVoter(tx.From(), tx.GetTxBody().ToAddress())*/
		}
	}
	// Add 1 to the number of blocks at this address
	/*dpos.dposStorage.SetTermWinnerMintCnt(block.Term, block.Signer)
	nextBlockTerm := (block.Time + blockInterval) / termInterval
	if nextBlockTerm == block.Term+1 {
		term := &Term{dPosStorage: dpos.dposStorage}
		term.kickOutValidator(block.Term)
	}*/
}

func (dpos *DPos) loadConfirmedBlockHeader(chain consensus.IChain) (*types.Header, error) {
	key, err := dpos.dposStorage.GetConfirmedBlockHash()
	if err != nil {
		return nil, err
	}
	header, err := chain.GetHeaderByHash(key)
	if err != nil {
		return nil, err
	}
	return header, nil
}

func (dpos *DPos) verifyBlockSigner(winner hasharry.Address, header *types.Header) error {
	if !types.VerifySigner(param.Net, winner, header.SignScript.PubKey) {
		return errors.New("not the signature of the address")
	}
	if !types.Verify(header.Hash, header.SignScript) {
		return errors.New("verify seal failed")
	}
	return nil
}

func (dpos *DPos) checkTime(lastHeader *types.Header, header *types.Header) error {
	nextSlot := NextSlot(header.Time)
	if lastHeader.Time >= nextSlot {
		return errors.New("create the future block")
	}
	if nextSlot-header.Time >= 1 {
		return fmt.Errorf("wait for last block arrived, next slot = %d, block time = %d ", nextSlot, header.Time)
	}
	if header.Time == nextSlot {
		return nil
	}
	return fmt.Errorf("wait for last block arrived, next slot = %d, block time = %d ", nextSlot, header.Time)
}

func (dpos *DPos) isSkipCurrentWinner(now uint64, lastTime uint64) bool {
	skipTimes := (now - lastTime) / param.SkipCurrentWinnerWaitTimeBase
	if skipTimes < 1 {
		return false
	}
	skipIndex := skipTimes % param.MaxWinnerSize
	return now%(param.SkipCurrentWinnerWaitTimeBase+skipIndex*param.BlockInterval) == 0
}

func (dpos *DPos) lookupWinners(now uint64) (hasharry.Address, error) {
	offset := now % param.TermInterval
	if offset%param.BlockInterval != 0 {
		return hasharry.Address{}, errors.New("invalid time to mint the block")
	}
	offset /= param.BlockInterval
	winners, err := dpos.dposStorage.GetTermWinners(now / param.TermInterval)
	if err != nil {
		return hasharry.Address{}, err
	}
	if len(winners.Candidates) == 0 {
		return hasharry.Address{}, errors.New("no winner to be found in storage")
	}
	offset %= uint64(len(winners.Candidates))
	winner := winners.Candidates[offset]
	return winner.Signer, nil
}

func (dpos *DPos) lookupWinnerNoExistToCreate(now uint64, parent *types.Header, chainReader consensus.IChain) (hasharry.Address, error) {
	offset := now % param.TermInterval
	if offset%param.BlockInterval != 0 {
		return hasharry.Address{}, errors.New("invalid time to mint the block")
	}
	offset /= param.BlockInterval
	winners, err := dpos.setWinners(now, parent, chainReader)
	if err != nil {
		return hasharry.Address{}, err
	}
	if len(winners) == 0 {
		return hasharry.Address{}, errors.New("no winner to be found in storage")
	}
	offset %= uint64(len(winners))
	winner := winners[offset]
	return winner.Signer, nil
}

// Get the hash of the last block of the previous cycle
// as the random number seed of the new cycle.
func (dpos *DPos) getPreTermLastBlockHash(current *types.Header, chain consensus.IChain) (*types.Header, error) {
	preTermLastHash, err := chain.GetTermLastHash(current.Term - 1)
	if err == nil {
		header, _ := chain.GetHeaderByHash(preTermLastHash)
		tHeader, _ := chain.GetHeaderByHeight(header.Height)
		if header.Height < current.Height && header.Hash.IsEqual(tHeader.Hash) {
			return header, nil
		}
	}

	// If the last block header of the last cycle cannot be obtained directly
	// from the chain, then look forward from the current block
	genesis, err := chain.GetHeaderByHeight(0)
	if err != nil {
		return nil, err
	}
	header, err := chain.GetHeaderByHeight(1)
	if err != nil {
		return genesis, nil
	}

	if header.Term >= current.Term {
		return genesis, nil
	}
	height := current.Height
	for height > 0 {
		height--
		header, err := chain.GetHeaderByHeight(height)
		if err != nil {
			continue
		}
		if header.Term < current.Term {
			return header, nil
		}
	}
	return nil, errors.New("not found")
}

func PrevSlot(now uint64) uint64 {
	return uint64((now-1)/param.BlockInterval) * param.BlockInterval
}

func NextSlot(now uint64) uint64 {
	return uint64((now+param.BlockInterval-1)/param.BlockInterval) * param.BlockInterval
}
