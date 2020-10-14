package types

// rlp encoded block for storage and transmission
type RlpBlock struct {
	*Header
	*RlpBody
}

func (rb *RlpBlock) TranslateToBlock() *Block {
	return &Block{
		Header: rb.Header,
		Body:   rb.RlpBody.TranslateToBody(),
	}
}

type RlpBlocks []*RlpBlock

func (rbs RlpBlocks) TranslateToBlocks() []*Block {
	blocks := []*Block{}
	for _, rlpBlock := range rbs {
		blocks = append(blocks, rlpBlock.TranslateToBlock())
	}
	return blocks
}
