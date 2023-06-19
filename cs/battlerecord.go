package cs

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

// BattleRecord is a recording of a single battle.
type BattleRecord struct {
	Num             int                         `json:"num,omitempty"`
	Tokens          []BattleRecordToken         `json:"tokens,omitempty"`
	ActionsPerRound [][]BattleRecordTokenAction `json:"actionsPerRound,omitempty"`
	tokensByNum     map[int]*BattleRecordToken
}

// A token on a battle board
type BattleRecordToken struct {
	Num              int       `json:"num,omitempty"`
	PlayerNum        int       `json:"playerNum,omitempty"`
	Token            ShipToken `json:"token,omitempty"`
	StartingPosition Vector    `json:"startingPosition,omitempty"`
}

// BattleRecordTokenAction represents an action for a token in a battle.
type BattleRecordTokenAction struct {
	Type              BattleRecordTokenActionType `json:"type,omitempty"`
	TokenNum          int                         `json:"tokenNum,omitempty"`
	From              Vector                      `json:"from,omitempty"`
	To                Vector                      `json:"to,omitempty"`
	Slot              int                         `json:"slot,omitempty"`
	TargetNum         int                         `json:"targetNum,omitempty"`
	TokensDestroyed   int                         `json:"tokensDestroyed,omitempty"`
	DamageDoneShields int                         `json:"damageDoneShields,omitempty"`
	DamageDoneArmor   int                         `json:"damageDoneArmor,omitempty"`
	TorpedoHits       int                         `json:"torpedoHits,omitempty"`
	TorpedoMisses     int                         `json:"torpedoMisses,omitempty"`
}

type BattleRecordTokenActionType int

const (
	TokenActionFire BattleRecordTokenActionType = iota
	TokenActionBeamFire
	TokenActionTorpedoFire
	TokenActionMove
	TokenActionRanAway
)

// SetupRecord populates a lookup table of items by guid.
func newBattleRecord(tokens []BattleRecordToken) *BattleRecord {
	tokensByNum := make(map[int]*BattleRecordToken)
	for i := range tokens {
		token := &tokens[i]
		tokensByNum[token.Num] = token
	}
	// always start with one round
	actionsPerRound := make([][]BattleRecordTokenAction, 1)
	actionsPerRound[0] = make([]BattleRecordTokenAction, 0)
	return &BattleRecord{
		Tokens:          tokens,
		ActionsPerRound: actionsPerRound,
		tokensByNum:     tokensByNum,
	}
}

func (a *BattleRecordTokenAction) String() string {
	return fmt.Sprintf("Action token %d from: %v to: %v target: %d", a.TokenNum, a.From, a.To, a.TargetNum)
}

// Add a new round to the record
func (b *BattleRecord) RecordNewRound() {
	b.ActionsPerRound = append(b.ActionsPerRound, []BattleRecordTokenAction{})
}

// set starting position of a token
func (b *BattleRecord) RecordStartingPosition(num int, position Vector) {
	b.tokensByNum[num].StartingPosition = position
}

// Record a move
func (b *BattleRecord) RecordMove(round int, token *battleToken, from, to Vector) {
	actions := b.ActionsPerRound[len(b.ActionsPerRound)-1]
	actions = append(actions, BattleRecordTokenAction{Type: TokenActionMove, TokenNum: token.Num, From: from, To: to})
	b.ActionsPerRound[len(b.ActionsPerRound)-1] = actions

	log.Debug().Msgf("Round: %d %v", round, b.ActionsPerRound[round][len(b.ActionsPerRound[round])-1])
}

// Record a token running away
func (b *BattleRecord) RecordRunAway(round int, token *battleToken) {
	actions := b.ActionsPerRound[len(b.ActionsPerRound)-1]
	actions = append(actions, BattleRecordTokenAction{Type: TokenActionRanAway, TokenNum: token.Num})
	b.ActionsPerRound[len(b.ActionsPerRound)-1] = actions

	log.Debug().Msgf("Round: %d %v", round, b.ActionsPerRound[round][len(b.ActionsPerRound[round])-1])

}

// Record a token firing a beam weapon
func (b *BattleRecord) RecordBeamFire(round int, token *battleToken, from Vector, to Vector, slot int, target *battleToken, damageDoneShields int, damageDoneArmor int, tokensDestroyed int) {
	actions := b.ActionsPerRound[len(b.ActionsPerRound)-1]
	actions = append(actions, BattleRecordTokenAction{Type: TokenActionBeamFire, TokenNum: token.Num, TargetNum: target.Num, From: from, To: to, Slot: slot, DamageDoneShields: damageDoneShields, DamageDoneArmor: damageDoneArmor, TokensDestroyed: tokensDestroyed})
	b.ActionsPerRound[len(b.ActionsPerRound)-1] = actions

	log.Debug().Msgf("Round: %d %v", round, b.ActionsPerRound[round][len(b.ActionsPerRound[round])-1])

}

// Record a token firing a salvo of torpedos
func (b *BattleRecord) RecordTorpedoFire(round int, token *battleToken, from Vector, to Vector, slot int, target *battleToken, damageDoneShields int, damageDoneArmor int, tokensDestroyed int, hits int, misses int) {
	actions := b.ActionsPerRound[len(b.ActionsPerRound)-1]
	actions = append(actions,
		BattleRecordTokenAction{
			Type:              TokenActionTorpedoFire,
			TokenNum:          token.Num,
			From:              from,
			To:                to,
			Slot:              slot,
			TargetNum:         target.Num,
			DamageDoneShields: damageDoneShields,
			DamageDoneArmor:   damageDoneArmor,
			TokensDestroyed:   tokensDestroyed,
			TorpedoHits:       hits,
			TorpedoMisses:     misses,
		},
	)
	b.ActionsPerRound[len(b.ActionsPerRound)-1] = actions

	log.Debug().Msgf("Round: %d %v", round, b.ActionsPerRound[round][len(b.ActionsPerRound[round])-1])

}
