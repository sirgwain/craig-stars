package cs

import "fmt"

type battleTokenAttribute int

const (
	battleTokenAttributeUnarmed       battleTokenAttribute = 0
	battleTokenAttributeArmed         battleTokenAttribute = 1 << 0
	battleTokenAttributeBomber        battleTokenAttribute = 1 << 1
	battleTokenAttributeFreighter     battleTokenAttribute = 1 << 2
	battleTokenAttributeStarbase      battleTokenAttribute = 1 << 3
	battleTokenAttributeFuelTransport battleTokenAttribute = 1 << 4
)

// a token for a battle
type battleToken struct {
	BattleRecordToken
	*ShipToken
	player            *Player
	designName        string // for String()
	cost              Cost
	attributes        battleTokenAttribute
	weaponSlots       []*battleWeaponSlot
	fleet             *Fleet
	attackPlayers     map[int]bool
	movesLeft         int
	movementMass      float64
	quantityDestroyed int
	destroyed         bool
	ranAway           bool
	movesMade         int
	armor             int
	shields           int
	stackShields      int
	totalStackShields int
	torpedoJamming    float64
	beamDefense       float64
}

// newBattleToken creates a new battle token from a shipToken. Mass and movement
// depend on cargo, so they are set when the battle is prepared.
func newBattleToken(rules *Rules, num int, position Vector, token *ShipToken, battlePlan BattlePlan, player *Player) *battleToken {
	battleToken := battleToken{
		BattleRecordToken: BattleRecordToken{
			Num:                     num,
			PlayerNum:               player.Num,
			Position:                position,
			DesignNum:               token.DesignNum,
			Initiative:              min(63, token.design.Spec.Initiative),
			Armor:                   token.design.Spec.Armor,
			StackShields:            token.design.Spec.Shields * token.Quantity,
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
		attributes:        token.battleAttributes(),
	}
	if !battleToken.hasWeapons() {
		battleToken.Tactic = BattleTacticDisengage
		battleToken.PrimaryTarget = BattleTargetNone
	}
	if battleToken.attributes&battleTokenAttributeStarbase != 0 {
		battleToken.Tactic = BattleTacticMaximizeDamage
		battleToken.PrimaryTarget = BattleTargetAny
		battleToken.SecondaryTarget = BattleTargetAny
		if !battleToken.hasWeapons() {
			battleToken.PrimaryTarget = BattleTargetNone
		}
	}

	// get the weapon slots for a token
	hull := rules.techs.GetHull(token.design.Hull)
	for _, slot := range token.design.Spec.WeaponSlots {
		weapon := rules.techs.GetHullComponent(slot.HullComponent)
		battleToken.weaponSlots = append(battleToken.weaponSlots, newBattleWeaponSlot(&battleToken, slot, weapon, hull.RangeBonus, token.design.Spec.TorpedoBonus, token.design.Spec.BeamBonus))
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

// battleAttributes returns the battle target classes this token's ships fall in
func (token *ShipToken) battleAttributes() battleTokenAttribute {
	spec := token.design.Spec
	attributes := getBattleTokenAttributes(spec.HullType, spec.HasWeapons)

	// Target classes describe the installed equipment, with armed ships taking priority.
	if !spec.HasWeapons {
		attributes &^= battleTokenAttributeBomber | battleTokenAttributeFreighter
		if len(spec.Bombs)+len(spec.SmartBombs)+len(spec.RetroBombs) > 0 {
			attributes |= battleTokenAttributeBomber
		} else if attributes&battleTokenAttributeFuelTransport == 0 && spec.CargoCapacity > 0 {
			attributes |= battleTokenAttributeFreighter
		}
	}
	return attributes
}

// isTargetOf returns true if any ship in this fleet would be a target of a battle plan's target
func (fleet *Fleet) isTargetOf(target BattleTarget) bool {
	for i := range fleet.Tokens {
		token := battleToken{attributes: fleet.Tokens[i].battleAttributes()}
		if token.isTargetOf(target) {
			return true
		}
	}
	return false
}

// convert hulltype to BattleTokenAttributes
func getBattleTokenAttributes(hullType TechHullType, hasWeapons bool) battleTokenAttribute {
	attributes := battleTokenAttributeUnarmed

	if hullType == TechHullTypeStarbase || hullType == TechHullTypeOrbitalFort {
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

// check if this token is still in the battle
func (token *battleToken) isStillInBattle() bool {
	return !token.destroyed && !token.ranAway
}

// beamDamageMultiplier returns the share of beam damage this token takes after
// deflectors. Designs without deflectors have a beamDefense of 0.
func (token *battleToken) beamDamageMultiplier() float64 {
	if token.beamDefense == 0 {
		return 1
	}
	return token.beamDefense
}

func (token *battleToken) getDistanceAway(position Vector) int {
	return max(Abs(token.Position.X-position.X), Abs(token.Position.Y-position.Y))
}

func (token *battleToken) String() string {
	return fmt.Sprintf("Player: %d, Token: %d %sx%d", token.PlayerNum, token.Num, token.designName, token.Quantity)
}

// willAttack reports whether this token's player is hostile toward another player
// in this battle, from attack orders, retaliation, and allied support.
func (token *battleToken) willAttack(otherPlayerNum int) bool {
	return otherPlayerNum != token.PlayerNum && token.attackPlayers[otherPlayerNum]
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
		return !token.hasWeapons() && ((token.attributes&battleTokenAttributeBomber) > 0 || (token.attributes&battleTokenAttributeFreighter) > 0)
	case BattleTargetUnarmedShips:
		return !token.hasWeapons() && token.attributes&battleTokenAttributeBomber == 0
	case BattleTargetFuelTransports:
		return !token.hasWeapons() && (token.attributes&battleTokenAttributeFuelTransport) > 0 && token.attributes&battleTokenAttributeBomber == 0
	case BattleTargetFreighters:
		return !token.hasWeapons() && (token.attributes&battleTokenAttributeFreighter) > 0 && token.attributes&(battleTokenAttributeBomber|battleTokenAttributeFuelTransport) == 0
	}

	return false
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
