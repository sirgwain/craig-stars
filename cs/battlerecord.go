package cs

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

// BattleRecord is a recording of a single battle.
type BattleRecord struct {
	Num             int                         `json:"num,omitempty"`
	PlanetNum       int                         `json:"planetNum,omitempty"`
	Position        Vector                      `json:"position,omitempty"`
	Tokens          []BattleRecordToken         `json:"tokens,omitempty"`
	ActionsPerRound [][]BattleRecordTokenAction `json:"actionsPerRound,omitempty"`
	Stats           struct {
		TokensDestroyedByPlayer map[int]int `json:"tokensDestroyedByPlayer,omitempty"`
	} `json:"stats,omitempty"`
}

// A token on a battle board
type BattleRecordToken struct {
	Num             int             `json:"num,omitempty"`
	PlayerNum       int             `json:"playerNum,omitempty"`
	DesignNum       int             `json:"designNum,omitempty"`
	Position        Vector          `json:"position,omitempty"`
	Initiative      int             `json:"initiative,omitempty"`
	Movement        int             `json:"movement,omitempty"`
	Tactic          BattleTactic    `json:"tactic,omitempty"`
	PrimaryTarget   BattleTarget    `json:"primaryTarget,omitempty"`
	SecondaryTarget BattleTarget    `json:"secondaryTarget,omitempty"`
	AttackWho       BattleAttackWho `json:"attackWho,omitempty"`
}

// BattleRecordTokenAction represents an action for a token in a battle.
type BattleRecordTokenAction struct {
	Type              BattleRecordTokenActionType `json:"type,omitempty"`
	TokenNum          int                         `json:"tokenNum,omitempty"`
	Round             int                         `json:"round,omitempty"`
	From              Vector                      `json:"from,omitempty"`
	To                Vector                      `json:"to,omitempty"`
	Slot              int                         `json:"slot,omitempty"`
	TargetNum         int                         `json:"targetNum,omitempty"`
	Target            *ShipToken                  `json:"target,omitempty"`
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

func (t BattleRecordTokenActionType) String() string {
	switch t {
	case TokenActionFire:
		return "Fire"
	case TokenActionBeamFire:
		return "BeamFire"
	case TokenActionTorpedoFire:
		return "TorpedoFire"
	case TokenActionMove:
		return "Move"
	case TokenActionRanAway:
		return "RanAway"
	default:
		return fmt.Sprintf("Unknown BattleRecordTokenActionType (%d)", t)
	}
}

// SetupRecord populates a lookup table of items by guid.
func newBattleRecord(num int, planetNum int, position Vector, tokens []BattleRecordToken) *BattleRecord {
	tokensByNum := make(map[int]*BattleRecordToken)
	for i := range tokens {
		token := &tokens[i]
		tokensByNum[token.Num] = token
	}
	// always start with one round
	actionsPerRound := make([][]BattleRecordTokenAction, 1)
	actionsPerRound[0] = make([]BattleRecordTokenAction, 0)
	return &BattleRecord{
		Num:             num,
		PlanetNum:       planetNum,
		Position:        position,
		Tokens:          tokens,
		ActionsPerRound: actionsPerRound,
	}
}

func (a BattleRecordTokenAction) String() string {
	return fmt.Sprintf("%s token %d from: %v to: %v target: %d", a.Type, a.TokenNum, a.From, a.To, a.TargetNum)
}

// Add a new round to the record
func (b *BattleRecord) RecordNewRound() {
	b.ActionsPerRound = append(b.ActionsPerRound, []BattleRecordTokenAction{})
}

// Record a move
func (b *BattleRecord) RecordMove(round int, token *battleToken, from, to Vector) {
	action := BattleRecordTokenAction{Type: TokenActionMove, Round: round, TokenNum: token.Num, From: from, To: to}
	actions := b.ActionsPerRound[len(b.ActionsPerRound)-1]
	actions = append(actions, action)
	b.ActionsPerRound[len(b.ActionsPerRound)-1] = actions

	log.Debug().Msgf("Round: %d %s", round, action)
}

// Record a token running away
func (b *BattleRecord) RecordRunAway(round int, token *battleToken) {
	action := BattleRecordTokenAction{Type: TokenActionRanAway, Round: round, TokenNum: token.Num}
	actions := b.ActionsPerRound[len(b.ActionsPerRound)-1]
	actions = append(actions, action)
	b.ActionsPerRound[len(b.ActionsPerRound)-1] = actions

	log.Debug().Msgf("Round: %d %s", round, action)

}

// Record a token firing a beam weapon
func (b *BattleRecord) RecordBeamFire(round int, token *battleToken, from Vector, to Vector, slot int, target battleToken, damageDoneShields int, damageDoneArmor int, tokensDestroyed int) {
	// copy the ship token into the record
	shipToken := *target.ShipToken

	action := BattleRecordTokenAction{Type: TokenActionBeamFire, Round: round, TokenNum: token.Num, TargetNum: target.Num, Target: &shipToken, From: from, To: to, Slot: slot, DamageDoneShields: damageDoneShields, DamageDoneArmor: damageDoneArmor, TokensDestroyed: tokensDestroyed}
	actions := b.ActionsPerRound[len(b.ActionsPerRound)-1]
	actions = append(actions, action)
	b.ActionsPerRound[len(b.ActionsPerRound)-1] = actions

	log.Debug().Msgf("Round: %d %s", round, action)

}

// Record a token firing a salvo of torpedos
func (b *BattleRecord) RecordTorpedoFire(round int, token *battleToken, from Vector, to Vector, slot int, target *battleToken, damageDoneShields int, damageDoneArmor int, tokensDestroyed int, hits int, misses int) {
	// copy the ship token into the record
	shipToken := *target.ShipToken
	action := BattleRecordTokenAction{
		Type:              TokenActionTorpedoFire,
		Round:             round,
		TokenNum:          token.Num,
		From:              from,
		To:                to,
		Slot:              slot,
		TargetNum:         target.Num,
		Target:            &shipToken,
		DamageDoneShields: damageDoneShields,
		DamageDoneArmor:   damageDoneArmor,
		TokensDestroyed:   tokensDestroyed,
		TorpedoHits:       hits,
		TorpedoMisses:     misses,
	}
	actions := b.ActionsPerRound[len(b.ActionsPerRound)-1]
	actions = append(actions, action)
	b.ActionsPerRound[len(b.ActionsPerRound)-1] = actions

	log.Debug().Msgf("Round: %d %s", round, action)

}
