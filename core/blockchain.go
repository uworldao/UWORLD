package core

import (
	"errors"
	"fmt"
	"github.com/uworldao/UWORLD/common/hasharry"
	"github.com/uworldao/UWORLD/consensus"
	"github.com/uworldao/UWORLD/core/types"
	"github.com/uworldao/UWORLD/database/blcdb"
	log "github.com/uworldao/UWORLD/log/log15"
	"github.com/uworldao/UWORLD/param"
	"sync"
)

const blockChainStorage = "block"

// Manage block chain
type BlockChain struct {
	currentHeight uint64
	stateRoot     hasharry.Hash
	contractRoot  hasharry.Hash
	consensusRoot hasharry.Hash
	accountState  IAccountState
	contractState IContractState
	consensus     consensus.IConsensus
	storage       IBlockChainStorage
	mutex         sync.RWMutex
	stateUpdateCh chan struct{}

	// List of transactions that have been stored and need
	// to be deleted by tx pool
	removeTxsCh chan types.Transactions

	// Confirmed valid block height
	confirmedHeight uint64
}

func NewBlockChain(dataDir string, consensus consensus.IConsensus, stateUpdateCh chan struct{},
	removeTxsCh chan types.Transactions, accountState IAccountState, contractState IContractState) (*BlockChain, error) {
	blockChain := &BlockChain{}
	storage := blcdb.NewBlockChainStorage(dataDir + "/" + blockChainStorage)
	err := storage.Open()
	if err != nil {
		return nil, err
	}

	blockChain.storage = storage
	blockChain.accountState = accountState
	blockChain.contractState = contractState
	blockChain.stateUpdateCh = stateUpdateCh
	blockChain.consensus = consensus
	blockChain.removeTxsCh = removeTxsCh
	stateRoot, _ := blockChain.storage.GetStateRoot()
	err = blockChain.accountState.InitTrie(stateRoot)
	if err != nil {
		return nil, err
	}
	blockChain.stateRoot = blockChain.accountState.RootHash()

	contractRoot, _ := blockChain.storage.GetContractRoot()
	err = blockChain.contractState.InitTrie(contractRoot)
	if err != nil {
		return nil, err
	}
	blockChain.contractRoot = blockChain.contractState.RootHash()

	consensusRoot, _ := blockChain.storage.GetConsensusRoot()
	err = blockChain.consensus.InitTrie(consensusRoot)
	if err != nil {
		return nil, err
	}
	blockChain.consensusRoot = blockChain.consensus.RootHash()

	if blockChain.currentHeight, err = blockChain.storage.GetLastHeight(); err != nil {
		if err := blockChain.SaveGenesisBlock(consensus.GetGenesisBlock()); err != nil {
			return nil, err
		}
	}

	blockChain.UpdateConfirmedHeight(consensus.GetConfirmedBlockHeader(blockChain).Height)

	return blockChain, nil
}

func (blc *BlockChain) CurrentHeader() (*types.Header, error) {
	blc.mutex.RLock()
	defer blc.mutex.RUnlock()

	return blc.storage.GetHeaderByHeight(blc.currentHeight)
}

func (blc *BlockChain) GetConfirmedHeight() uint64 {
	blc.mutex.RLock()
	defer blc.mutex.RUnlock()

	return blc.confirmedHeight
}

func (blc *BlockChain) GetHistoryConfirmedHeight(height uint64) (uint64, error) {
	return blc.storage.GetHistoryConfirmedHeight(height)
}

func (blc *BlockChain) UpdateConfirmedHeight(height uint64) {
	blc.mutex.Lock()
	defer blc.mutex.Unlock()

	blc.confirmedHeight = height
	blc.accountState.UpdateConfirmedHeight(height)
}

func (blc *BlockChain) GetLastHeight() uint64 {
	blc.mutex.RLock()
	defer blc.mutex.RUnlock()

	return blc.currentHeight
}

func (blc *BlockChain) GetHeaderByHeight(height uint64) (*types.Header, error) {
	if height > blc.GetLastHeight() {
		return nil, errors.New("not exist")
	}
	return blc.storage.GetHeaderByHeight(height)
}

func (blc *BlockChain) GetHeaderByHash(hash hasharry.Hash) (*types.Header, error) {
	return blc.storage.GetHeader(hash)
}

func (blc *BlockChain) GetBlockByHeight(height uint64) (*types.Block, error) {
	header, err := blc.GetHeaderByHeight(height)
	if err != nil {
		return nil, err
	}
	txs, err := blc.storage.GetTransactions(header.TxRoot)
	if err != nil {
		return nil, err
	}
	rlpBody := &types.RlpBody{txs}
	block := &types.Block{Header: header, Body: rlpBody.TranslateToBody()}
	return block, nil
}

func (blc *BlockChain) GetBlockByHash(hash hasharry.Hash) (*types.Block, error) {
	header, err := blc.GetHeaderByHash(hash)
	if err != nil {
		return nil, err
	}
	txs, err := blc.storage.GetTransactions(header.TxRoot)
	if err != nil {
		return nil, err
	}
	rlpBody := &types.RlpBody{txs}
	block := &types.Block{Header: header, Body: rlpBody.TranslateToBody()}
	return block, nil
}

func (blc *BlockChain) GetTransaction(hash hasharry.Hash) (types.ITransaction, error) {
	rlpTx, err := blc.storage.GetTransaction(hash)
	if err != nil {
		return nil, err
	}
	return rlpTx.TranslateToTransaction(), nil
}

func (blc *BlockChain) GetTransactionIndex(hash hasharry.Hash) (types.ITransactionIndex, error) {
	txIndex, err := blc.storage.GetTxLocation(hash)
	if err != nil {
		return nil, err
	}
	if txIndex.Height > blc.GetLastHeight() {
		return nil, errors.New("not exist")
	}
	return txIndex, nil
}

func (blc *BlockChain) GetRlpBlockByHeight(height uint64) (*types.RlpBlock, error) {
	header, err := blc.storage.GetHeaderByHeight(height)
	if err != nil {
		return nil, err
	}
	txs, err := blc.storage.GetTransactions(header.TxRoot)
	if err != nil {
		return nil, err
	}
	rlpBody := &types.RlpBody{txs}
	block := &types.RlpBlock{Header: header, RlpBody: rlpBody}
	return block, nil
}

func (blc *BlockChain) GetRlpBlockByHash(hash hasharry.Hash) (*types.RlpBlock, error) {
	header, err := blc.storage.GetHeader(hash)
	if err != nil {
		return nil, err
	}
	txs, err := blc.storage.GetTransactions(header.TxRoot)
	if err != nil {
		return nil, err
	}
	return &types.RlpBlock{Header: header, RlpBody: &types.RlpBody{txs}}, nil
}

func (blc *BlockChain) GetAddressVote(address hasharry.Address) uint64 {
	var vote uint64
	state := blc.accountState.GetAccountState(address)
	vote += state.GetBalance(param.Token.String())
	return vote
}

func (blc *BlockChain) GetTermLastHash(term uint64) (hasharry.Hash, error) {
	return blc.storage.GetTermLastHash(term)
}

func (blc *BlockChain) InsertChain(block *types.Block) error {
	if _, err := blc.GetBlockByHeight(block.Height - 1); err != nil {
		return err
	}
	return blc.dealBlock(block)
}

func (blc *BlockChain) saveBlock(block *types.Block) error {
	blc.mutex.Lock()
	defer blc.mutex.Unlock()

	blc.storage.UpdateHeader(block.Header)
	blc.storage.UpdateTransactions(block.TxRoot, block.Body.TranslateToRlpBody().Transactions)
	blc.storage.UpdateTxLocation(block.GetTxsLocations())
	blc.storage.UpdateHeightHash(block.Height, block.Hash)
	blc.storage.UpdateHistoryConfirmedHeight(block.Height, blc.confirmedHeight)
	blc.storage.UpdateTermLastHash(block.Term, block.Hash)
	blc.stateRoot, _ = blc.accountState.StateTrieCommit()
	blc.contractRoot, _ = blc.contractState.ContractTrieCommit()
	blc.consensusRoot, _ = blc.consensus.Commit()
	blc.storage.UpdateStateRoot(blc.stateRoot)
	blc.storage.UpdateContractRoot(blc.contractRoot)
	blc.storage.UpdateConsensusRoot(blc.consensusRoot)

	blc.currentHeight = block.Height
	blc.storage.UpdateLastHeight(block.Height)

	log.Info("Save block", "height", block.Height, "hash", block.HashString(),
		"state root", block.StateRoot.String(),
		"contract root", block.ContractRoot.String(),
		"consensus root", block.ConsensusRoot.String(),
		"signer", block.Signer.String(), "txcount", block.Transactions.Len(),
		"time", block.Time, "term", block.Term)
	return nil
}

func (blc *BlockChain) SaveGenesisBlock(block *types.Block) error {
	blc.mutex.Lock()
	defer blc.mutex.Unlock()

	if err := blc.VerifyGenesis(block); err != nil {
		return err
	}
	blc.updateGenesisState(block)
	blc.storage.UpdateHeader(block.Header)
	blc.storage.UpdateTransactions(block.TxRoot, block.Body.TranslateToRlpBody().Transactions)
	blc.storage.UpdateTxLocation(block.GetTxsLocations())
	blc.storage.UpdateHeightHash(block.Height, block.Hash)
	blc.storage.UpdateLastHeight(block.Height)
	blc.storage.UpdateHistoryConfirmedHeight(block.Height, 0)
	blc.consensus.SetConfirmedHeader(block.Header)
	blc.stateRoot, _ = blc.accountState.StateTrieCommit()
	blc.contractRoot, _ = blc.contractState.ContractTrieCommit()
	blc.consensusRoot, _ = blc.consensus.Commit()

	blc.storage.UpdateConsensusRoot(blc.consensusRoot)
	blc.currentHeight = block.Height
	log.Info("Save block", "height", block.Height, "hash", block.HashString(),
		"state", block.StateRoot.String(), "signer", block.Signer.String(), "txcount", block.Transactions.Len(),
		"time", block.Time, "term", block.Term)
	return nil
}

func (blc *BlockChain) updateState(block *types.Block) error {
	for _, tx := range block.Body.Transactions {
		switch tx.GetTxType() {
		case types.NormalTransaction:
			if err := blc.accountState.UpdateFrom(tx, block.Height); err != nil {
				return err
			}
			if err := blc.accountState.UpdateTo(tx, block.Height); err != nil {
				return err
			}
		case types.ContractTransaction:
			if err := blc.accountState.UpdateFrom(tx, block.Height); err != nil {
				return err
			}
			if err := blc.accountState.UpdateTo(tx, block.Height); err != nil {
				return err
			}
			blc.contractState.UpdateContract(tx, block.Height)
			/*case types.VoteToCandidate:
				fallthrough
			case types.LoginCandidate:
				fallthrough
			case types.LogoutCandidate:
				if err := blc.accountState.UpdateFrom(tx, block.Height); err != nil {
					return err
				}*/
		}

	}
	return blc.accountState.UpdateFees(block.Body.Transactions.SumFees(), block.Height)
}

func (blc *BlockChain) updateGenesisState(block *types.Block) error {
	for _, tx := range block.Body.Transactions {
		switch tx.GetTxType() {
		case types.NormalTransaction:
			if err := blc.accountState.UpdateTo(tx, block.Height); err != nil {
				return err
			}
		}
	}
	blc.consensus.UpdateConsensus(block)
	return nil
}

func (blc *BlockChain) updateConsensus(block *types.Block) error {
	blc.consensus.UpdateConsensus(block)
	return nil
}

func (blc *BlockChain) UpdateCurrentBlockHeight(height uint64) {
	blc.mutex.Lock()
	defer blc.mutex.Unlock()

	blc.currentHeight = height
	blc.storage.UpdateLastHeight(height)
}

func (blc *BlockChain) ValidationBlockHash(header *types.Header) (bool, error) {
	localHeader, err := blc.GetHeaderByHeight(header.Height)
	if err != nil {
		return false, err
	}
	return localHeader.Hash.IsEqual(header.Hash), nil
}

func (blc *BlockChain) VerifyGenesis(block *types.Block) error {
	var sumCoins uint64
	for _, tx := range block.Transactions {
		sumCoins += tx.GetTxBody().GetAmount()
	}
	if sumCoins != param.GenesisCoins {
		return fmt.Errorf("wrong genesis coins")
	}
	return nil
}

func (blc *BlockChain) verifyBlock(block *types.Block) error {
	if block.Height <= blc.GetLastHeight() {
		return ErrDuplicateBlock
	}
	if !block.VerifyTxRoot() {
		log.Warn("tx root wrong", "height", block.Header.Height, "tx root", block.Header.StateRoot.String())
		return errors.New("wrong tx root")
	}
	if !block.StateRoot.IsEqual(blc.StateRoot()) {
		log.Warn("state root wrong", "height", block.Header.Height, "state root", block.Header.StateRoot.String())
		//blc.accountState.Print()
		return errors.New("wrong state root")
	}
	if !block.ContractRoot.IsEqual(blc.ContractRoot()) {
		log.Warn("contract root wrong", "height", block.Header.Height, "contract root", block.Header.ContractRoot.String())
		return errors.New("wrong contract root")
	}
	if !block.ConsensusRoot.IsEqual(blc.ConsensusRoot()) {
		log.Warn("consensus root wrong", "height", block.Header.Height, "consensus root", block.Header.ConsensusRoot.String())
		return errors.New("wrong consensus root")
	}
	if err := blc.verifyTxs(block.Transactions, block.Height); err != nil {
		return err
	}
	parent, err := blc.GetHeaderByHash(block.ParentHash)
	if err != nil {
		return ErrNoParent
	}
	if err = blc.consensus.VerifyHeader(block.Header, parent); err != nil {
		return err
	}
	if err = blc.consensus.VerifySeal(blc, block.Header, parent); err != nil {
		return err
	}
	return nil
}

func (blc *BlockChain) verifyTx(tx types.ITransaction, blockHeight uint64) error {
	if err := tx.VerifyTx(); err != nil {
		return err
	}

	if err := blc.consensus.VerifyTx(tx); err != nil {
		return err
	}

	if err := blc.accountState.VerifyState(tx); err != nil {
		return err
	}

	if err := blc.contractState.VerifyState(tx); err != nil {
		return err
	}

	if err := blc.verifyBusiness(tx, blockHeight); err != nil {
		return err
	}

	return nil
}

func (blc *BlockChain) verifyBusiness(tx types.ITransaction, blockHeight uint64) error {
	switch tx.GetTxType() {
	case types.NormalTransaction:
		account := blc.accountState.GetAccountState(tx.From())
		return account.VerifyNonce(tx.GetNonce())
	}
	return nil
}

func (blc *BlockChain) verifyTxs(txs types.Transactions, blockHeight uint64) error {
	address := make(map[string]bool)
	for _, tx := range txs {
		if tx.IsCoinBase() {
			if err := blc.verifyCoinBaseTx(tx, blockHeight, 0); err != nil {
				return err
			}
		} else {
			if err := blc.verifyTx(tx, blockHeight); err != nil {
				blc.removeTxsCh <- types.Transactions{tx}
				return err
			}
		}
		from := tx.From().String()
		if _, ok := address[from]; !ok {
			address[from] = true
		} else {
			return errors.New("one address in a block can only send one transaction")
		}
	}
	return nil
}

func (blc *BlockChain) verifyCoinBaseTx(tx types.ITransaction, height, sumFees uint64) error {
	return tx.VerifyCoinBaseTx(height, sumFees)
}

// When a serious inconsistency occurs, it can fall back to any height
// below the effective height
func (blc *BlockChain) FallBackTo(height uint64) error {
	confirmedHeight := blc.confirmedHeight
	if height > confirmedHeight && height != 0 {
		err := fmt.Sprintf("the height of the fallback must be less than or equal %d and greater than %d", confirmedHeight, 0)
		log.Error("Fall back to block height", "height", height, "error", err)
		return errors.New(err)
	}

	var curBlockHeight, nextBlockHeight uint64
	curStateRoot := hasharry.Hash{}
	curContractRoot := hasharry.Hash{}
	curConsensusRoot := hasharry.Hash{}

	nextBlockHeight = height + 1
	curBlockHeight = height

	// set new confirmed height and header
	hisConfirmedHeight, err := blc.GetHistoryConfirmedHeight(curBlockHeight)
	if err != nil {
		log.Error("Fall back to block height", "height", height, "error", "can not find history confirmed height")
		return fmt.Errorf("fall back to block height %d failed! Can not find history confirmed height", height)
	}
	hisHeader, err := blc.GetHeaderByHeight(hisConfirmedHeight)
	if err != nil {
		log.Error("Fall back to block height", "height", height, "error", "can not find block")
		return fmt.Errorf("fall back to block height %d failed! Can not find block %d", height, hisConfirmedHeight)
	}
	blc.consensus.SetConfirmedHeader(hisHeader)

	log.Warn("Fall back to block height", "height", height)
	header, err := blc.GetHeaderByHeight(nextBlockHeight)
	if err != nil {
		log.Error("Fall back to block height", "height", height, "error", "can not find block")
		return fmt.Errorf("fall back to block height %d failed! Can not find block %d", height, nextBlockHeight)
	}

	blc.mutex.Lock()
	defer blc.mutex.Unlock()

	blc.confirmedHeight = hisConfirmedHeight
	blc.accountState.UpdateConfirmedHeight(hisConfirmedHeight)

	// fall back to pre state root
	curStateRoot = header.StateRoot
	err = blc.accountState.InitTrie(curStateRoot)
	if err != nil {
		log.Error("Fall back to block height", "height", height, "error", "init state trie failed")
		return fmt.Errorf("fall back to block height %d failed! nit state trie failed", height)
	}
	blc.stateRoot = blc.accountState.RootHash()

	// fall back to pre contract root
	curContractRoot = header.ContractRoot
	err = blc.contractState.InitTrie(curContractRoot)
	if err != nil {
		log.Error("Fall back to block height", "height", height, "error", "init contract trie failed")
		return fmt.Errorf("fall back to block height %d failed! nit contract trie failed", height)
	}
	blc.contractRoot = blc.contractState.RootHash()

	// fall back to consensus root
	curConsensusRoot = header.ConsensusRoot
	err = blc.consensus.InitTrie(curConsensusRoot)
	if err != nil {
		log.Error("Fall back to block height", "height", height, "error", "init consensus trie failed")
		return fmt.Errorf("fall back to block height %d failed! init consensus trie failed", height)
	}
	blc.consensusRoot = blc.consensus.RootHash()

	blc.currentHeight = curBlockHeight
	blc.storage.UpdateLastHeight(curBlockHeight)
	return nil
}

func (blc *BlockChain) FallBack() {
	blc.FallBackTo(blc.GetConfirmedHeight())
}

func (blc *BlockChain) dealBlock(block *types.Block) error {
	err := blc.verifyBlock(block)
	if err == nil {
		if err := blc.updateState(block); err != nil {
			return err
		}
		blc.updateConsensus(block)
		blc.saveBlock(block)
		blc.stateUpdateCh <- struct{}{}
		return nil
	}
	return err
}

func (blc *BlockChain) StateRoot() hasharry.Hash {
	blc.mutex.RLock()
	defer blc.mutex.RUnlock()

	return blc.stateRoot
}

func (blc *BlockChain) ContractRoot() hasharry.Hash {
	blc.mutex.RLock()
	defer blc.mutex.RUnlock()

	return blc.contractRoot
}

func (blc *BlockChain) ConsensusRoot() hasharry.Hash {
	blc.mutex.RLock()
	defer blc.mutex.RUnlock()

	return blc.consensusRoot
}

func (blc *BlockChain) TireRoot() (hasharry.Hash, hasharry.Hash, hasharry.Hash) {
	blc.mutex.RLock()
	defer blc.mutex.RUnlock()

	return blc.stateRoot, blc.contractRoot, blc.consensusRoot
}

func (blc *BlockChain) CloseStorage() error {
	blc.mutex.RLock()
	defer blc.mutex.RUnlock()

	var err error
	err = blc.contractState.Close()
	err = blc.accountState.Close()
	err = blc.consensus.Close()
	err = blc.storage.Close()
	log.Info("Close storage")
	return err
}
