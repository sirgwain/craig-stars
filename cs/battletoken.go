package cs

import (
	"fmt"
	"math"
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

// newBattleToken creates a new battle token from a shipToken.
func newBattleToken(num int, position BattleVector, cargoMass int, token *ShipToken, battlePlan BattlePlan, player *Player, techFinder TechFinder) *battleToken {
	battleToken := battleToken{
		BattleRecordToken: BattleRecordToken{
			Num:                     num,
			PlayerNum:               player.Num,
			Position:                position,
			DesignNum:               token.DesignNum,
			Initiative:              token.design.Spec.Initiative,
			Mass:                    token.design.Spec.Mass + cargoMass,
			Armor:                   token.design.Spec.Armor,
			StackShields:            token.design.Spec.Shields * token.Quantity,
			Movement:                token.design.getMovement(cargoMass),
			StartingQuantity:        token.Quantity,
			StartingQuantityDamaged: token.QuantityDamaged,
			StartingDamage:          int(token.Damage),
			Tactic:                  battlePlan.Tactic,
			PrimaryTarget:           battlePlan.PrimaryTarget,
			SecondaryTarget:         battlePlan.SecondaryTarget,
			AttackWho:               battlePlan.AttackWho,
		},
		ShipToken:         token,
		player:            player,
		designName:        token.design.Name,
		cost:              token.design.Spec.Cost,
		armor:             token.design.Spec.Armor,
		shields:           token.design.Spec.Shields,
		stackShields:      token.Quantity * token.design.Spec.Shields,
		totalStackShields: token.Quantity * token.design.Spec.Shields,
		torpedoJamming:    token.design.Spec.TorpedoJamming,
		beamDefense:       token.design.Spec.BeamDefense,
		attributes:        getBattleTokenAttributes(token.design.Spec.HullType, token.design.Spec.HasWeapons),
	}

	// get the weapon slots for a token
	weaponSlots := make([]*battleWeaponSlot, 0)
	hull := techFinder.GetHull(token.design.Hull)
	if len(token.design.Spec.WeaponSlots) > 0 {
		minRange := math.MaxInt
		maxRange := 0
		for _, slot := range token.design.Spec.WeaponSlots {
			weapon := techFinder.GetHullComponent(slot.HullComponent)
			bws := newBattleWeaponSlot(&battleToken, slot, weapon, hull.RangeBonus, token.design.Spec.TorpedoBonus, token.design.Spec.BeamBonus)
			weaponSlots = append(weaponSlots, bws)
			minRange = MinInt(minRange, bws.weaponRange)
			maxRange = MaxInt(maxRange, bws.weaponRange)
			if bws.weaponType == battleWeaponTypeBeam {
				battleToken.attributes |= battleTokenAttributeHasBeams
			} else if bws.weaponType == battleWeaponTypeTorpedo {
				battleToken.attributes |= battleTokenAttributeHasTorpedos
			}
		}
		battleToken.weaponSlots = weaponSlots
		battleToken.minRange = minRange
		battleToken.maxRange = maxRange

		// to maximize our damage, we either close in all the way
		// or get close enough so all our weapons can fire
		if battleToken.hasBeamWeapons() {
			battleToken.maxDamageRange = 0
		} else {
			battleToken.minRange = 0
		}
	}

	return &battleToken
}

// getCargoPerShip returns the amount of cargo a single ship in a stack is carrying
func getCargoPerShip(fleetCargo, fleetCargoCapacity, tokenCargoCapacity int) int {
	// add cargo from the fleet to each token
	cargoMass := 0
	if fleetCargo > 0 && tokenCargoCapacity > 0 {
		// see how much this ship's cargo capacity is compared to the fleet total
		shipCargoPercent := float64(tokenCargoCapacity) / float64(fleetCargoCapacity)
		cargoMass = int(float64(fleetCargo) * shipCargoPercent)
	}

	return cargoMass
}

// convert hulltype to BattleTokenAttributes
func getBattleTokenAttributes(hullType TechHullType, hasWeapons bool) battleTokenAttribute {
	attributes := battleTokenAttributeUnarmed

	if hullType == TechHullTypeStarbase {
		attributes |= battleTokenAttributeStarbase
	}

	if hasWeapons {
		attributes |= battleTokenAttributeArmed
	}

	if hullType == TechHullTypeFreighter {
		attributes |= battleTokenAttributeFreighter
	}

	if hullType == TechHullTypeFuelTransport {
		attributes |= battleTokenAttributeFuelTransport
	}

	if hullType == TechHullTypeBomber {
		attributes |= battleTokenAttributeBomber
	}

	return attributes
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
