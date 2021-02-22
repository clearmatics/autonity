package types

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/rlp"
	"io"
)

type ProofType uint64
const (
	InnocentProof ProofType = iota
	ChallengeProof
	AccusationProof
)

type Rule uint8
const (
	PN Rule = iota
	PO
	PVN
	PVO
	C
	C1

	GarbageMessage		// message was signed by valid member, but it cannot be decoded.
	InvalidProposal     // The value proposed by proposer cannot pass the blockchain's validation.
	InvalidProposer     // A proposal sent from none proposer nodes of the committee.
	Equivocation        // Multiple distinguish votes(proposal, prevote, precommit) sent by validator.
	UnknownRule
)

// The proof used by accountability precompiled contract to validate the proof of innocent or misbehavior.
// Since precompiled contract take raw bytes as input, so it should be rlp encoded into bytes before client send the
// proof into autonity contract.
type RawProof struct {
	Rule       Rule       // rule id.
	Message    []byte      // the raw rlp encoded msg to be considered as suspicious one
	Evidence   [][]byte    // the raw rlp encoded msgs as proof of innocent or misbehavior.
}

func (p *RawProof) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{p.Rule, p.Message, p.Evidence})
}

func (p *RawProof) DecodeRLP(s *rlp.Stream) error {
	var proof struct {
		Rule Rule
		Message []byte
		Evidence [][]byte
	}
	if err := s.Decode(&proof); err != nil {
		return err
	}

	p.Rule, p.Message, p.Evidence = proof.Rule, proof.Message, proof.Evidence
	return nil
}

// Proof is what to prove that one is misbehaving, one should be slashed when a valid proof is rise.
type Proof struct {
	Rule       Rule
	Message    ConsensusMessage   // the msg to be considered as suspicious one
	Evidence   []ConsensusMessage // the msgs as proof of innocent or misbehavior.
}

// OnChainProof to be stored by autonity contract for on-chain proof management.
type OnChainProof struct {
	Sender  common.Address `abi:"sender"`
	Msghash common.Hash    `abi:"msghash"`
	// rlp enoded bytes for struct Rawproof object.
	Rawproof []byte        `abi:"rawproof"`
}

// event to submit proofs via standard transaction.
type SubmitProofEvent struct {
	Proofs []OnChainProof
	Type ProofType
}