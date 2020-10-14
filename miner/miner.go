package miner

import (
	"github.com/jhdriver/UWORLD/common/hasharry"
	"github.com/jhdriver/UWORLD/consensus"
	"github.com/jhdriver/UWORLD/core"
	"github.com/jhdriver/UWORLD/core/types"
	"github.com/jhdriver/UWORLD/crypto/ecc/secp256k1"
	log "github.com/jhdriver/UWORLD/log/log15"
	"github.com/jhdriver/UWORLD/param"
	"time"
)

const resultChanSize = 10
const maxBlockTransactions = 999

// Generate block miner
type Miner struct {
	// Node private key
	key *secp256k1.PrivateKey
	// Node address
	signer      hasharry.Address
	consensus   consensus.IConsensus
	blockChain  core.IBlockChain
	txPool      core.ITxPool
	recvBlock   chan *types.Block
	genBlkCh    chan *types.Block
	minerWorkCh chan bool
	stop        chan bool
	isStop      chan bool
}

func NewMiner(consensus consensus.IConsensus, blockChain core.IBlockChain, txPool core.ITxPool, key *secp256k1.PrivateKey, signer hasharry.Address,
	genBlkCh chan *types.Block, minerWorkCh chan bool) *Miner {
	return &Miner{
		signer:      signer,
		key:         key,
		blockChain:  blockChain,
		txPool:      txPool,
		consensus:   consensus,
		recvBlock:   make(chan *types.Block, resultChanSize),
		genBlkCh:    genBlkCh,
		minerWorkCh: minerWorkCh,
		stop:        make(chan bool, 1),
		isStop:      make(chan bool, 2),
	}
}

func (miner *Miner) Start() error {
	//go miner.control()
	go miner.work()
	go miner.waitBlock()

	log.Info("Miner startup successful")
	return nil
}

func (miner *Miner) Stop() error {
	close(miner.stop)
	log.Info("Stop miner")
	return nil
}

// Check every second whether it is a block
func (miner *Miner) work() {
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case _, _ = <-miner.stop:
			miner.isStop <- true
			log.Info("Stop minting block")
			return
		case now := <-ticker:
			miner.mintingBlock(uint64(now.Unix()))
		}
	}
}

// Send the generated block
func (miner *Miner) waitBlock() {
	for {
		select {
		case _, _ = <-miner.stop:
			miner.isStop <- true
			log.Info("Stop wait block")
			return
		default:
			for block := range miner.recvBlock {
				log.Info("Create block", "height", block.Height, "hash", block.HashString(), "state", block.StateRoot.String(), "signer", block.Signer.String(), "txcount", block.Transactions.Len(),
					"time", block.Time, "term", block.Term)
				miner.genBlkCh <- block
			}
		}
	}
}

func (miner *Miner) mintingBlock(now uint64) {
	currentHeader, err := miner.blockChain.CurrentHeader()
	if err != nil {
		log.Error("Get current block failed!", "err", err)
		return
	}
	if currentHeader.Time > now {
		now = currentHeader.Time + 1
	}
	stateRoot, contractRoot, consensusRoot := miner.blockChain.TireRoot()

	// Build block header
	header := &types.Header{
		StateRoot:     stateRoot,
		ContractRoot:  contractRoot,
		ConsensusRoot: consensusRoot,
		ParentHash:    currentHeader.Hash,
		Height:        currentHeader.Height + 1,
		Time:          now,
		Term:          now / miner.consensus.GetTermInterval(),
		Signer:        miner.signer,
	}
	// Check if it is your turn to make blocks
	err = miner.consensus.CheckWinner(miner.blockChain, header)
	if err != nil {
		//.Warn("check winner failed!", "height", header.Height, "error", err)
		return
	}
	// Generate block
	if block, err := miner.generateBlock(header); err != nil {
		log.Error("Generate the block, failed!", "height", header.Height, "err", err)
		return
	} else {
		miner.recvBlock <- block
	}
}

func (miner *Miner) generateBlock(header *types.Header) (*types.Block, error) {
	txs := miner.getTransactions(header.Height)
	header.TxRoot = txs.Hash()
	header.SetHash()
	block := types.NewBlock(header, types.NewBody(txs))
	// Sign the generated block
	err := miner.consensus.Sign(block)
	return block, err
}

// Get transactions from the transaction pool and generate coinbase transactions
func (miner *Miner) getTransactions(height uint64) types.Transactions {
	txs := miner.txPool.Gets(maxBlockTransactions)
	coinBase := miner.getCoinBase(txs, height)
	coinBaseTx := miner.generateCoinBaseTx(coinBase)
	coinBaseTx.SetHash()
	txs = append(txs, coinBaseTx)
	return txs
}

func (miner *Miner) getCoinBase(txs types.Transactions, height uint64) uint64 {
	//return types.CalCoinBase(height) + txs.SumFees()
	return types.CalCoinBase(height, param.CoinHeight)
}

func (miner *Miner) generateCoinBaseTx(coinBase uint64) types.ITransaction {
	return &types.Transaction{
		TxHead: &types.TransactionHead{
			TxHash:     hasharry.Hash{},
			TxType:     types.NormalTransaction,
			From:       hasharry.StringToAddress(types.CoinBase),
			Nonce:      0,
			Fees:       0,
			Time:       uint64(time.Now().Unix()),
			Note:       "",
			SignScript: &types.SignScript{},
		},
		TxBody: &types.NormalTransactionBody{
			Contract: param.Token,
			To:       miner.signer,
			Amount:   coinBase,
		},
	}
}
