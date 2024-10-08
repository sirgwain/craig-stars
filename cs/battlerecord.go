package cs

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

// BattleRecord is a recording of a single battle.
type BattleRecord struct {
	Num             int                          `json:"num,omitempty"`
	PlanetNum       int                          `json:"planetNum,omitempty"`
	Position        Vector                       `json:"position,omitempty"`
	Tokens          []BattleRecordToken          `json:"tokens,omitempty"`
	ActionsPerRound [][]BattleRecordTokenAction  `json:"actionsPerRound,omitempty"`
	DestroyedTokens []BattleRecordDestroyedToken `json:"destroyedTokens,omitempty"`
	Stats           BattleRecordStats            `json:"stats,omitempty"`
}
type BattleRecordStats struct {
	NumPlayers             int           `json:"numPlayers,omitempty"`
	NumShipsByPlayer       map[int]int   `json:"numShipsByPlayer,omitempty"`
	ShipsDestroyedByPlayer map[int]int   `json:"shipsDestroyedByPlayer,omitempty"`
	DamageTakenByPlayer    map[int]int   `json:"damageTakenByPlayer,omitempty"`
	CargoLostByPlayer      map[int]Cargo `json:"cargoLostByPlayer,omitempty"`
}

// A token on a battle board
type BattleRecordToken struct {
	Num                     int             `json:"num,omitempty"`
	PlayerNum               int             `json:"playerNum,omitempty"`
	DesignNum               int             `json:"designNum,omitempty"`
	Position                BattleVector    `json:"position,omitempty"`
	Initiative              int             `json:"initiative,omitempty"`
	Mass                    int             `json:"mass,omitempty"`
	Armor                   int             `json:"armor,omitempty"`
	StackShields            int             `json:"stackShields,omitempty"`
	Movement                int             `json:"movement,omitempty"`
	StartingQuantity        int             `json:"startingQuantity,omitempty"`
	StartingQuantityDamaged int             `json:"startingQuantityDamaged,omitempty"`
	StartingDamage          int             `json:"startingDamage,omitempty"`
	Tactic                  BattleTactic    `json:"tactic,omitempty"`
	PrimaryTarget           BattleTarget    `json:"primaryTarget,omitempty"`
	SecondaryTarget         BattleTarget    `json:"secondaryTarget,omitempty"`
	AttackWho               BattleAttackWho `json:"attackWho,omitempty"`
}

type BattleRecordDestroyedToken struct {
	Num       int `json:"num,omitempty"`
	PlayerNum int `json:"playerNum,omitempty"`
	DesignNum int `json:"designNum,omitempty"`
	Quantity  int `json:"quantity,omitempty"`
	design    *ShipDesign
}

// BattleRecordTokenAction represents an action for a token in a battle.
type BattleRecordTokenAction struct {
	Type              BattleRecordTokenActionType `json:"type,omitempty"`
	TokenNum          int                         `json:"tokenNum,omitempty"`
	Round             int                         `json:"round,omitempty"`
	From              BattleVector                `json:"from,omitempty"`
	To                BattleVector                `json:"to,omitempty"`
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

type BattleVector struct {
	X int `json:"x"`
	Y int `json:"y"`
}

var BattleVectorRight BattleVector = BattleVector{1, 0}
var BattleVectorLeft BattleVector = BattleVector{-1, 0}
var BattleVectorUp BattleVector = BattleVector{0, 1}
var BattleVectorDown BattleVector = BattleVector{0, -1}
var BattleVectorUpRight BattleVector = BattleVector{1, 1}
var BattleVectorUpLeft BattleVector = BattleVector{-1, 1}
var BattleVectorDownRight BattleVector = BattleVector{1, -1}
var BattleVectorDownLeft BattleVector = BattleVector{-1, -1}

func (v1 BattleVector) Add(v2 BattleVector) BattleVector {
	return BattleVector{v1.X + v2.X, v1.Y + v2.Y}
}

func (v1 BattleVector) distance(v2 BattleVector) int {
	return MaxInt(AbsInt(v1.X-v2.X), AbsInt(v1.Y-v2.Y))
}

func (v BattleVector) scale(scale int) BattleVector {
	return BattleVector{v.X * scale, v.Y * scale}
}

// SetupRecord populates a lookup table of items by guid.
func newBattleRecord(num int, planetNum int, position Vector, tokens []BattleRecordToken) *BattleRecord {
	numShipsByPlayer := make(map[int]int, 2)
	tokensByNum := make(map[int]*BattleRecordToken)
	playerNums := make(map[int]bool, 2)
	for i := range tokens {
		token := &tokens[i]
		tokensByNum[token.Num] = token
		numShipsByPlayer[token.PlayerNum] = numShipsByPlayer[token.PlayerNum] + token.StartingQuantity
		playerNums[token.PlayerNum] = true
	}
	// always start with one round
	actionsPerRound := make([][]BattleRecordTokenAction, 1)
	actionsPerRound[0] = make([]BattleRecordTokenAction, 0)

	// setup some default stats
	shipsDestroyedByPlayer := make(map[int]int, len(playerNums))
	damageTakenByPlayer := make(map[int]int, len(playerNums))
	cargoLostByPlayer := make(map[int]Cargo, len(playerNums))

	for playerNum := range playerNums {
		shipsDestroyedByPlayer[playerNum] = 0
		damageTakenByPlayer[playerNum] = 0
		cargoLostByPlayer[playerNum] = Cargo{}
	}

	return &BattleRecord{
		Num:             num,
		PlanetNum:       planetNum,
		Position:        position,
		Tokens:          tokens,
		ActionsPerRound: actionsPerRound,
		Stats: BattleRecordStats{
			NumPlayers:             len(playerNums),
			NumShipsByPlayer:       numShipsByPlayer,
			ShipsDestroyedByPlayer: shipsDestroyedByPlayer,
			DamageTakenByPlayer:    damageTakenByPlayer,
			CargoLostByPlayer:      cargoLostByPlayer,
		},
	}
}

func (a BattleRecordTokenAction) String() string {
	return fmt.Sprintf("%s token %d from: %v to: %v target: %d", a.Type, a.TokenNum, a.From, a.To, a.TargetNum)
}

// Add a new round to the record
func (b *BattleRecord) recordNewRound() {
	b.ActionsPerRound = append(b.ActionsPerRound, []BattleRecordTokenAction{})
}

// Record a move
func (b *BattleRecord) recordMove(round int, token *battleToken, from, to BattleVector) BattleRecordTokenAction {
	action := BattleRecordTokenAction{Type: TokenActionMove, Round: round, TokenNum: token.Num, From: from, To: to}
	actions := b.ActionsPerRound[len(b.ActionsPerRound)-1]
	actions = append(actions, action)
	b.ActionsPerRound[len(b.ActionsPerRound)-1] = actions
	return action
}

// Record a token running away
func (b *BattleRecord) recordRunAway(round int, token *battleToken) BattleRecordTokenAction {
	action := BattleRecordTokenAction{Type: TokenActionRanAway, Round: round, TokenNum: token.Num}
	actions := b.ActionsPerRound[len(b.ActionsPerRound)-1]
	actions = append(actions, action)
	b.ActionsPerRound[len(b.ActionsPerRound)-1] = actions

	return action
}

// Record a token firing a beam weapon
func (b *BattleRecord) recordBeamFire(round int, token *battleToken, from BattleVector, to BattleVector, slot int, target battleToken, damageDoneShields int, damageDoneArmor int, tokensDestroyed int) {
	// copy the ship token into the record
	shipToken := *target.ShipToken

	action := BattleRecordTokenAction{Type: TokenActionBeamFire, Round: round, TokenNum: token.Num, TargetNum: target.Num, Target: &shipToken, From: from, To: to, Slot: slot, DamageDoneShields: damageDoneShields, DamageDoneArmor: damageDoneArmor, TokensDestroyed: tokensDestroyed}
	actions := b.ActionsPerRound[len(b.ActionsPerRound)-1]
	actions = append(actions, action)
	b.ActionsPerRound[len(b.ActionsPerRound)-1] = actions
	b.Stats.DamageTakenByPlayer[token.PlayerNum] += damageDoneArmor

	log.Debug().Msgf("Round: %d %s", round, action)

}

// Record a token firing a salvo of torpedos
func (b *BattleRecord) recordTorpedoFire(round int, token *battleToken, from BattleVector, to BattleVector, slot int, target *battleToken, damageDoneShields int, damageDoneArmor int, tokensDestroyed int, hits int, misses int) {
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
	b.Stats.DamageTakenByPlayer[token.PlayerNum] += damageDoneArmor

	log.Debug().Msgf("Round: %d %s", round, action)
}

// record a destroyed token for the battle message
func (b *BattleRecord) recordDestroyedToken(token *battleToken, quantityDestroyed int) {
	destroyedTokens := BattleRecordDestroyedToken{}
	destroyedTokens.Num = token.Num
	destroyedTokens.PlayerNum = token.PlayerNum
	destroyedTokens.DesignNum = token.design.Num
	destroyedTokens.Quantity = quantityDestroyed
	destroyedTokens.design = token.design
	b.DestroyedTokens = append(b.DestroyedTokens, destroyedTokens)
	b.Stats.ShipsDestroyedByPlayer[token.PlayerNum] += quantityDestroyed
}
