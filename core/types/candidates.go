package types

import (
	"github.com/jhdriver/UWORLD/common/hasharry"
)

// Super nodes
type Winners struct {
	Candidates      []*Candidate
	ElectParentHash hasharry.Hash
}

type Candidate struct {
	Signer hasharry.Address
	PeerId string
	Weight uint64
}

type Candidates struct {
	Members []*Candidate
}

func NewCandidates() *Candidates {
	return &Candidates{Members: make([]*Candidate, 0)}
}

func (c *Candidates) Set(newMem *Candidate) {
	for _, mem := range c.Members {
		if mem.Signer.IsEqual(newMem.Signer) {
			return
		}
	}
	c.Members = append(c.Members, newMem)
}

func (c *Candidates) Remove(reMem *Candidate) {
	for i, mem := range c.Members {
		if mem.Signer.IsEqual(reMem.Signer) {
			c.Members = append(c.Members[0:i], c.Members[i+1:]...)
			return
		}
	}
}

func (c *Candidates) Len() int {
	return len(c.Members)
}

type SortableCandidates []*Candidate

func (p SortableCandidates) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p SortableCandidates) Len() int      { return len(p) }
func (p SortableCandidates) Less(i, j int) bool {
	if p[i].Weight < p[j].Weight {
		return false
	} else if p[i].Weight > p[j].Weight {
		return true
	} else {
		return p[i].Weight < p[j].Weight
	}
}
