package afd

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
)

type MsgStore struct {
	// map[Height]map[Round]map[MsgType]map[common.address]*ConsensusMessage
	messages map[uint64]map[int64]map[uint64]map[common.Address]*types.ConsensusMessage
}

// store msg into msg store, it returns msg that is equivocation than the input msg, and an errEquivocation.
// otherwise it return nil, nil
func(ms *MsgStore) Save(m *types.ConsensusMessage) (*types.ConsensusMessage, error) {
	height, _ := m.Height()
	roundMap, ok := ms.messages[height.Uint64()]
	if !ok {
		roundMap = make(map[int64]map[uint64]map[common.Address]*types.ConsensusMessage)
		ms.messages[height.Uint64()] = roundMap
	}

	round, _ := m.Round()
	msgTypeMap, ok := roundMap[round]
	if !ok {
		msgTypeMap = make(map[uint64]map[common.Address]*types.ConsensusMessage)
		roundMap[round] = msgTypeMap
	}

	addressMap, ok := msgTypeMap[m.Code]
	if !ok {
		addressMap = make(map[common.Address]*types.ConsensusMessage)
		msgTypeMap[m.Code] = addressMap
	}

	msg, ok := addressMap[m.Address]
	if !ok {
		addressMap[m.Address] = m
		return nil, nil
	}

	// check equivocation here.
	if types.RLPHash(msg.Payload()) != types.RLPHash(m.Payload()) {
		return msg, errEquivocation
	}
	return nil, nil
}

func(ms *MsgStore) removeMsg(m *types.ConsensusMessage) {
	height, _ := m.Height()
	round, _ := m.Round()
	delete(ms.messages[height.Uint64()][round][m.Code], m.Address)
}

func(ms *MsgStore) DeleteMsgsAtHeight(height uint64) {
	// Remove all messgages at this height
	for _, msgTypeMap := range ms.messages[height] {
		for _, addressMap := range msgTypeMap {
			for _, m := range addressMap {
				ms.removeMsg(m)
			}
		}
	}
	// Delete map entry for this height
	delete(ms.messages, height)
}

// get take height and query conditions to query those msgs from msg store, it returns those msgs satisfied the condition.
func (ms *MsgStore) Get(height uint64, query func(*types.ConsensusMessage) bool) []types.ConsensusMessage {

	var result []types.ConsensusMessage
	roundMap, ok := ms.messages[height]
	if !ok {
		return result
	}

	for _, msgTypeMap := range roundMap {
		for _, addressMap := range msgTypeMap {
			for _, m := range addressMap {
				if query(m) {
					result = append(result, *m)
				}
			}
		}
	}

	return result
}