package cs

import (
	"fmt"
)

type battleTokenAttribute int

const (
	battleTokenAttributeUnarmed       battleTokenAttribute = 0
	battleTokenAttributeArmed         battleTokenAttribute = 1 << 0
	battleTokenAttributeBomber        battleTokenAttribute = 1 << 1
	battleTokenAttributeFreighter     battleTokenAttribute = 1 << 2
	battleTokenAttributeStarbase      battleTokenAttribute = 1 << 3
	battleTokenAttributeFuelTransport battleTokenAttribute = 1 << 4
	battleTokenAttributeHasBeams      battleTokenAttribute = 1 << 5
	battleTokenAttributeHasTorpedos   battleTokenAttribute = 1 << 6
)

// a token for a battle
type battleToken struct {
	BattleRecordToken
	*ShipToken
	player            *Player
	designName        string // for String()
	cost              Cost
	attributes        battleTokenAttribute
	moveTarget        *battleToken
	weaponSlots       []*battleWeaponSlot
	quantityDestroyed int
	damaged           bool
	destroyed         bool
	ranAway           bool
	movesMade         int
	armor             int
	shields           int
	stackShields      int
	totalStackShields int
	minRange          int
	maxRange          int
	maxDamageRange    int
	torpedoJamming    float64
	beamDefense       float64
}

func (token *battleToken) hasWeapons() bool {
	return (token.attributes & battleTokenAttributeArmed) > 0
}

func (token *battleToken) hasBeamWeapons() bool {
	return (token.attributes & battleTokenAttributeHasBeams) > 0
}

// check if this token is still in the battle
func (token *battleToken) isStillInBattle() bool {
	return !token.destroyed && !token.ranAway
}

func (token *battleToken) getDistanceAway(position BattleVector) int {
	return MaxInt(AbsInt(token.Position.X-position.X), AbsInt(token.Position.Y-position.Y))
}

func (token *battleToken) String() string {
	return fmt.Sprintf("Player: %d, Token: %d %sx%d", token.PlayerNum, token.Num, token.designName, token.Quantity)
}

// return true if this fleet will attack a fleet by another player based on player
// relations and the fleet battle plan
func (token *battleToken) willAttack(otherPlayerNum int) bool {
	// if we have weapons and we don't own this other fleet, see if we
	// would target it
	player := token.player
	if token.hasWeapons() && token.Tactic != BattleTacticDisengage && otherPlayerNum != player.Num {
		switch token.AttackWho {
		case BattleAttackWhoEnemies:
			return player.IsEnemy(otherPlayerNum)
		case BattleAttackWhoEnemiesAndNeutrals:
			return player.IsEnemy(otherPlayerNum) || player.IsNeutral(otherPlayerNum)
		case BattleAttackWhoEveryone:
			return true
		}
	}
	return false
}

// isTargetOf returns true if the BattleOrder Target type would target this token
func (token *battleToken) isTargetOf(target BattleTarget) bool {
	switch target {
	case BattleTargetAny:
		return true
	case BattleTargetNone:
		return false
	case BattleTargetStarbase:
		return (token.attributes & battleTokenAttributeStarbase) > 0
	case BattleTargetArmedShips:
		return (token.attributes & battleTokenAttributeArmed) > 0
	case BattleTargetBombersFreighters:
		return (token.attributes&battleTokenAttributeBomber) > 0 || (token.attributes&battleTokenAttributeFreighter) > 0
	case BattleTargetUnarmedShips:
		return (token.attributes & battleTokenAttributeArmed) == 0
	case BattleTargetFuelTransports:
		return (token.attributes & battleTokenAttributeFuelTransport) > 0
	case BattleTargetFreighters:
		return (token.attributes & battleTokenAttributeFreighter) > 0
	}

	return false
}

// willTarget returns true if this token will target the target token
func (token *battleToken) willTarget(target *battleToken) bool {
	return token.willAttack(target.PlayerNum) && (target.isTargetOf(token.PrimaryTarget) || target.isTargetOf(token.SecondaryTarget))
}

// find all weapon targets for this token and return the best one
func (token *battleToken) findWeaponsTargets(tokens []*battleToken) {
	for _, weapon := range token.weaponSlots {
		targets := weapon.findTargets(tokens)
		weapon.targets = targets
	}
}

// find the most attractive target for this token to move towards
func (token *battleToken) findMoveTarget() *battleToken {

	var bestTarget *battleToken
	var bestAttractiveness float64
	for _, weapon := range token.weaponSlots {
		if len(weapon.targets) == 0 {
			continue
		}
		bestWeaponTarget := weapon.targets[0]
		attractiveness := weapon.getAttractiveness(bestWeaponTarget)
		if bestTarget == nil || attractiveness > bestAttractiveness {
			bestAttractiveness = attractiveness
			bestTarget = bestWeaponTarget
		}
	}

	return bestTarget
}

// regenerateShields regenerates the shields of the given token if the player regenerates shields
// and the token has shields.
func (token *battleToken) regenerateShields() {
	player := token.player

	if player.Race.Spec.ShieldRegenerationRate > 0 && token.stackShields > 0 {
		regenerationAmount := int(float64(token.totalStackShields)*player.Race.Spec.ShieldRegenerationRate + 0.5)
		token.stackShields = int(Clamp(token.stackShields+regenerationAmount, 0, token.totalStackShields))
	}
}
