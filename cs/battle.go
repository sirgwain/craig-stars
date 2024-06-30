package cs

import (
	"fmt"
	"math"
	"sort"

	"github.com/rs/zerolog/log"
)

// From: https://wiki.starsautohost.org/wiki/Guts_of_the_Battle_Engine
// ===================================================================
// Here are the guts of the battle engine as I understand it from both experience, observation and the help file
// (please pull me up on any points I get wrong)
//
// For a battle to take place 2 or more fleets (or a fleet and a starbase) must be at the same location and at
// least one of the fleets must be armed and have orders to attack ships of the others race (the type of ships
// involved doesn't matter). If are race has a fleet present at a location where there is a battle, but doesn't
// have orders to attack any of the other races there and none of the other races present has orders to attack it
// then it will not take part in the battle (and can not benefit from potential tech gain -- actually you can benefit
// from tech gain, a fact I learned from trying not to get the tech gain in a wolf/lamb tech exchange - LEit).
// Each ship present at the battle will form part of a token (AKA a stack), it is possible to have a token comprised
// of just a single ship. Tokens are always of ships of the same design. Each ship design in each fleet will create
// a token, splitting a few ships off to form a second fleet before the battle will create a second token on the
// battle board.
//
// The battle grid is made up of 10 squares by 10 squares. Each token is in a single square, there can be more
// than one token in the same square.There is an limit of 256 tokens per battle event for all players involved,
// if this limit is exceeded, then excess tokens will be left out (those created from fleets with the highest
// fleet numbers), in such a case each player will have an equal number of tokens, each player will be guaranteed
// to get their "share" of the available token slots (ie in a 4 race battle 256 / 4 = 64 token slots), if a race
// doesn't use up all their "slots" then they are shared equally between the other players.
//
// Each battle is made up of rounds. There are a maximum of 16 rounds in each battle. Each round has two parts,
// movement and shooting. Each token has a speed rating, and will be able to move between 0 and 3 squares in a
// single turn. If a token has a fractional speed rating then they will get a bonus square of movement every set
// number of turns. a 1/4 bonus means an extra square of movement on the first round and then on every fourth round
// after that starting with the fifth. A 1/2 speed bonus gets a bonus square of movement every other turn starting
// with the first, and a 3/4 speed bonus gets a bonus square of movement for the first three rounds of every 4 round
// cycle. The order of movement is this, each token with 3 movement squares moves a single square, then each token
// with 2+ movement moves a single square (if it had speed 3 then it would move for its second square) and then all
// ships with at least one square of movement move again. At each stage the ships with the most weight will move first
// though there is less than a 15% difference in weight then there is a chance that the lighter ship will go first.
// The smaller the weight % difference the greater the chance of the lighter ship going first.
//
// Each token has an attractiveness rating. This is used in both working out where ships move to and which ships are
// shot at first. The essence of the formula is cost / defence. A ship will have different attractiveness ratings
// verses different types of weapons (beams, sappers, torpedoes and capital missiles). Cost is calculated by summing
// the resource and boranium costs of the ship design used (iron and germ costs don't affect the attractiveness rating).
// Defence is calculated by the shield and armour dp modified by the enemies torpedo accuracy (after base accuracy,
// comps and jammers are worked out) if defending vs torps or capital missiles, the effects of double damage for unshielded
// targets vs capital missiles and the effects of deflectors against beam weapons. The attractiveness rating can be
// change during the course of the battle as shields and armour deplete. Attractiveness doesn't take into account the
// one missile one kill rule, thus chaff has become a fairly effective tactic.
//
// battle orders are comprised of 4 parts. A primary and secondary target type, legitimate races to attack and the tactic
// to use in battle. Ships will only attack tokens belonging to legitimate target races, however if another race present
// has any ships (including unarmed ships) with battles orders to attack your race then that race will also be considered
// a legitimate target. When attacking ships will try and shoot the most attractive ship of a type listed as a primary
// target and if no ships are available which are primary targets then the most attractive ship of a type listed as a
// secondary target will be targeted. Ships which are not listed as primary or secondary targets will not get shot at,
// even if they are shooting back.
//
// There are 6 different battle orders which determine the movement AI of the ships in battle, the movement AI is applied
// each time a ship wants to move a square on the battle board.:
//   - Disengage - If there is any enemy ship in firing range then move to any square further away than your current square.
//     If you are in range of an enemy weapon but cannot move further away then try move to a square that is of
//     the same distance away. If you are in range of the enemies weapons and cannot move away or maintain distance
//     then move to a random square. If you are not in range of the enemies weapons then move randomly. Also you
//     will try and disengage which will require 7 squares of movement to be clocked up before you can leave from
//     the battle board.
//   - Disengage if Challenged - Behaves like Maximise Damage until token takes damage and then behaves like Disengage.
//   - Minimise Damage to Self - (Not 100% sure on this one) If within range of an enemy weapon then move away from the
//     enemy (just like Disengage). If out of range of the enemies weapons or cannot move away from
//     the enemy then try and get in range of the best available target without moving towards the enemy.
//   - Maximise Net Damage - Locate most attractive primary target (or secondary if no primary targets are left). If out
//     of range with ANY weapon then move towards target. If in range with all weapons them move as to
//     maximise damage_done/damage_taken. The effect of this is if your weapons are longer range then
//     try to stay at maximum range. If your weapons range is the same then do random movement while
//     staying in range. If your weapons are shorter range and also beam weapons then attempt to close
//     in to zero range.
//   - Maximise Damage Ratio - As Maximise Net Damage but only considers the longest range weapon.
//   - Maximise Damage - Locate most attractive primary target (or secondary if no primary targets are left). If any of
//     your weapons are out of range of that token then keep moving to squares that are closer to it until
//     in range with all weapons. If using any beam weapons (as they have range dissipation) then attempt
//     to close to 0 range. If just using missiles or torps and in range then move randomly to a squares
//     still in range.
//
// Note that there is a bug when fighting starbases, the battle AI doesn't count the +1 range bonus when calculating movement.
// This mainly applies when your ships are attempting to get out of range of the enemy, so vs starbase with range 6 missiles,
// your ships will move to distance 7, the movement AI won't calculate that they are still in range even when they keep getting
// shot at.
// After the movement phase all ships will shoot their weapons, a token will fire all weapons from the same slot in a single
// shot. The weapon slot with the highest initiative will fire first. If there are two ships with slots of the same init,
// then the ships will be randomly given a priority over who can fire first (which will stick for the entire battle). The
// rest of the weapon slots are then fired in init order. Damage is worked out in between each shot and applied to the ships.
// If ships or tokens are destroyed before their turn to shoot then they won't be able to fire back. The movement AI will
// go after the most attractive primary target on the board, but if this token is not in range, then the ship will fire on
// the most attractive primary target within range (or secondary if none available). Starbases have a +1 range bonus to all
// their weapons (this also gets applied to minefield sweeping rates), though cannot move. The movement AI doesn't take this
// bonus into account when moving ships to close in on an enemy starbase.
//
// Damage for each shot is calculated by multiplying the number of weapons in the slot by the number of ships in the token
// by the amount of dp the weapon does. For beam weapons, this damage will dissipate by 10% over the range of the beam
// (for a range 2 beam - no dissipation at range 0, 5% dissipation at range 1 and 10% dissipation at range 2). Also
// capacitors and deflectors will modify the damage actually done to the enemy ship. Damage will be applied first to the
// tokens shield stack and then to armour only when the entire shield stack of the token is down. For missile ships,
// each missile fired will be tested to see if it will hit, the chance to hit is based on the base accuracy, the computers
// on the ship and the enemy jammers. Missiles that miss will do 1/8 of their damage to the shields and won't affect armour.
// For missiles that hit, upto half will be taken by the shields, the rest will go to the armour. For capital missiles
// any damage done after the shields are taken down will do double damage to the armour. Whole ship kills are worked out
// by adding up all the damage done to the armour by a single salvo (from a token's slot) and dividing this by the amount
// of armour each single ship in the token has left (total armour x token damage %). The number of complete ships the shot
// could kill will be removed from the enemy token, the rest of the damage will divided equally among the rest of the ships
// in the token and applied as damage. As token armour is stored in 1/512ths (about 0.2%s) of total armour and not as an
// exact dp figure (shields are stored as an exact figure), there may be some rounding of the damage after each salvo
// (AFAIK its always rounds up). This fact can be abused by creating lots of small fleet tokens with weak missiles
// and many slots, where each slot that hits will do 0.2% damage to the enemy token even if each individual missile
// would do less damage normally (especially the case with a beta torp shooting a large nub stack).
//
// After all the weapons that are in range have fired, the next round begins, starting with ship movement.
// The battle is ended when either the 16 round timer runs out, there is only one race left present on the battle board
// or if there are two or more races which have no hostile intentions towards each other.
//
// After a battle, salvage is created. This is equal to 1/3 of the current mineral costs of all the ships that where
// destroyed during the battle. This is left at the location of the battle and will decay over time, or if the battle
// happened over a planet, then the minerals will get deposited there.
//
// Any races that took part in a battle and had at least one ship that managed to survive (either through surviving
// till the end or retreating beforehand) has a potential to gain tech levels from ships that where destroyed during
// the battle. For the exact details of the formulas and chances involved see the Guts of Tech Trading.
//
// Movement speed and moves per round.
// 3/4 is      1110
// 1 is        1111
// 1 1/4 is    2111
// 1 1/2 is    2121
// 1 3/4 is    2212
// 2 is        2222
// 2 1/4 is    3222
// 2 1/2 is    3232
//
// ===================================================================

const battleWidth, battleHeight = 10, 10

// The battler interface is the main entrypoint into running battles
// A batter is created for fleets in the same location in the universe and it is used
// to check if a battle will occur, and run the battle.
// Running a battle returns a BattleRecord that is passed along to each player.
type battler interface {
	hasTargets() bool
	runBattle() *BattleRecord
}

// battle defines the state of a battle as it progresses
type battle struct {
	planet   *Planet
	position Vector
	tokens   []*battleToken
	round    int
	players  map[int]*Player
	board    [battleWidth][battleHeight]int // the number of tokens in each square
	record   *BattleRecord
	rules    *Rules
}

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
	attributes        battleTokenAttribute
	moveTarget        *battleToken
	targetedBy        []*battleToken
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

func (bt *battleToken) hasWeapons() bool {
	return (bt.attributes & battleTokenAttributeArmed) > 0
}

func (bt *battleToken) hasBeamWeapons() bool {
	return (bt.attributes & battleTokenAttributeHasBeams) > 0
}

// check if this token is still in the battle
func (bt *battleToken) isStillInBattle() bool {
	return !bt.destroyed && !bt.ranAway
}

func (bt *battleToken) getDistanceAway(position BattleVector) int {
	return MaxInt(AbsInt(bt.Position.X-position.X), AbsInt(bt.Position.Y-position.Y))
}

func (bt *battleToken) String() string {
	return fmt.Sprintf("Player: %d %sx%d", bt.PlayerNum, bt.design.Name, bt.Quantity)
}

type battleWeaponType int

const (
	battleWeaponTypeBeam battleWeaponType = iota
	battleWeaponTypeTorpedo
)

// A token firing weapons
type battleWeaponSlot struct {
	// The token with the weapon
	token *battleToken

	// The weapon slot
	slot ShipDesignSlot

	// how many of this weapon are in this slot
	quantity int

	// The type of weapon this weapon slot is
	weaponType battleWeaponType

	// The range of this weapon
	weaponRange int

	// the power of the weapon
	power int

	// true if this weapon damages shields only (i.e. a sapper)
	damagesShieldsOnly bool

	// the accuracy of the weapon, if it's a torpedo
	accuracy float64

	// the initiative of the weapon
	initiative int

	// gattling guns hit all targets in range
	hitsAllTargets bool

	// capital ships missiles do double damage after shields are gone
	capitalShipMissile bool
}

// Return true if this weapon slot wiil damage this token
// if this is a sapper and the target is out of shields, return false
func (slot *battleWeaponSlot) willDamage(target *battleToken) bool {
	if target == nil {
		return false
	}
	return !slot.damagesShieldsOnly || (slot.damagesShieldsOnly && target.stackShields > 0)
}

// Return true if this weapon slot is in range of the token target
func (slot *battleWeaponSlot) isInRange(target *battleToken) bool {
	if target == nil {
		return false
	}
	return slot.isInRangePosition(target.Position)
}

func (slot *battleWeaponSlot) isInRangePosition(position BattleVector) bool {
	// diagonal shots count as one move, so we take the max distance on the x or y as our actual distance away
	// i.e. 4 over, 1 up is 4 range away, 3 over 2 up is 3 range away, etc.
	return slot.token.getDistanceAway(position) <= slot.weaponRange
}

func (slot *battleWeaponSlot) isInRangeValue(rangeValue int) bool {
	return rangeValue <= slot.weaponRange
}

var positionsByPlayer = []BattleVector{
	{1, 4},
	{8, 5},
	{4, 1},
	{5, 8},
	{1, 1},
	{8, 1},
	{8, 8},
	{1, 8},
}

var movementByRound = [9][4]int{
	{1, 0, 1, 0},
	{1, 1, 0, 1},
	{1, 1, 1, 1},
	{2, 1, 1, 1},
	{2, 1, 2, 1},
	{2, 2, 1, 2},
	{2, 2, 2, 2},
	{3, 2, 2, 2},
	{3, 2, 3, 2},
}

// get the movement of this design with additional cargo
func getBattleMovement(idealEngineSpeed, movementBonus, mass, numEngines int) int {
	return Clamp(((idealEngineSpeed+movementBonus)-2)-((mass)/numEngines/70), 2, 10)
}

// BuildBattle builds a battle recording with all the battle tokens for a list of fleets that contains more than one player.
// We'll use this to determine if a battle should take place at this location.
// Also, any players that have a potential battle will discover each other's designs.
func newBattler(rules *Rules, techFinder TechFinder, battleNum int, players map[int]*Player, fleets []*Fleet, planet *Planet) battler {
	if len(fleets) == 0 {
		log.Error().Msg("Can't build battle with no fleets.")
		return nil
	}

	// track
	sortedPlayerNums := make([]int, 0, len(players))
	for _, player := range players {
		sortedPlayerNums = append(sortedPlayerNums, player.Num)
	}
	sort.Ints(sortedPlayerNums)

	playerStartingPositions := make(map[int]BattleVector)
	for i, num := range sortedPlayerNums {
		if i >= len(positionsByPlayer) {
			log.Warn().Msg("Oh noes! We have a battle with more players than we have positions for...")
		}
		playerStartingPositions[num] = positionsByPlayer[i%len(positionsByPlayer)]
	}

	board := [battleWidth][battleHeight]int{}
	// add each fleet's token to the battle
	tokens := []*battleToken{}
	tokenRecords := []BattleRecordToken{}
	num := 0
	dampening := 0
	for _, fleet := range fleets {
		totalCargo := fleet.Cargo.Total()
		totalCargoCapacity := fleet.Spec.CargoCapacity

		for i := range fleet.Tokens {
			token := &fleet.Tokens[i]
			num++

			cargoMass := 0
			if totalCargo > 0 && token.design.Spec.CargoCapacity > 0 {
				// see how much this ship's cargo capacity is compared to the fleet total
				shipCargoPercent := float64(token.design.Spec.CargoCapacity) / float64(totalCargoCapacity)
				cargoMass = int(float64(totalCargo) * shipCargoPercent)
			}

			position := playerStartingPositions[fleet.PlayerNum]
			battleToken := &battleToken{
				BattleRecordToken: BattleRecordToken{
					Num:              num,
					PlayerNum:        fleet.PlayerNum,
					Position:         position,
					DesignNum:        token.DesignNum,
					Initiative:       token.design.Spec.Initiative,
					Mass:             token.design.Spec.Mass + cargoMass,
					Movement:         token.design.getMovement(cargoMass),
					StartingQuantity: token.Quantity,
					Tactic:           fleet.battlePlan.Tactic,
					PrimaryTarget:    fleet.battlePlan.PrimaryTarget,
					SecondaryTarget:  fleet.battlePlan.SecondaryTarget,
					AttackWho:        fleet.battlePlan.AttackWho,
				},
				ShipToken:         token,
				armor:             token.design.Spec.Armor,
				shields:           token.design.Spec.Shields,
				stackShields:      token.Quantity * token.design.Spec.Shields,
				totalStackShields: token.Quantity * token.design.Spec.Shields,
				torpedoJamming:    token.design.Spec.TorpedoJamming,
				beamDefense:       token.design.Spec.BeamDefense,
				attributes:        getBattleTokenAttributes(token.design.Spec.HullType, token.design.Spec.HasWeapons),
			}
			board[position.X][position.Y] += battleToken.StartingQuantity
			// get the weapon slots for a token
			weaponSlots := make([]*battleWeaponSlot, 0)
			hull := techFinder.GetHull(token.design.Hull)
			if len(token.design.Spec.WeaponSlots) > 0 {
				minRange := math.MaxInt
				maxRange := 0
				var lastWeapon *battleWeaponSlot
				for slotIndex, slot := range token.design.Spec.WeaponSlots {
					weapon := techFinder.GetHullComponent(slot.HullComponent)
					if slotIndex != 0 && lastWeapon != nil && weapon != nil && lastWeapon.slot.HullComponent == weapon.Name {
						// this weapon is the same, just add it's quantity
						lastWeapon.quantity += slot.Quantity
						continue
					}
					bws := newBattleWeaponSlot(battleToken, slot, techFinder.GetHullComponent(slot.HullComponent), hull.RangeBonus, token.design.Spec.TorpedoInaccuracyFactor, token.design.Spec.BeamBonus)
					lastWeapon = bws
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

			// find the highest dampener we have
			dampening = MaxInt(dampening, token.design.Spec.ReduceMovement)

			tokens = append(tokens, battleToken)
			tokenRecords = append(tokenRecords, battleToken.BattleRecordToken)
		}
	}

	// apply dampening
	if dampening > 0 {
		for _, token := range tokens {
			// we only dampen movement of ships that move, not starbases (obviously)
			// and we can't go below 2
			if token.Movement > 0 {
				token.Movement = Clamp(token.Movement-dampening, 2, 10)
			}
		}
	}

	planetNum := 0
	if planet != nil {
		planetNum = planet.Num
	}

	return &battle{
		planet:   planet,
		position: fleets[0].Position,
		tokens:   tokens,
		record:   newBattleRecord(battleNum, planetNum, fleets[0].Position, tokenRecords),
		players:  players,
		rules:    rules,
	}
}

// newBattleWeaponSlot creates a new BattleWeaponSlot object
func newBattleWeaponSlot(token *battleToken, slot ShipDesignSlot, hc *TechHullComponent, rangeBonus int, torpedoInaccuracyFactor float64, beamBonus float64) *battleWeaponSlot {
	weaponSlot := &battleWeaponSlot{
		token:              token,
		slot:               slot,
		weaponRange:        hc.Range + rangeBonus,
		power:              hc.Power,
		damagesShieldsOnly: hc.DamageShieldsOnly,
		accuracy:           (100.0 - (100.0-float64(hc.Accuracy))*torpedoInaccuracyFactor) / 100.0,
		initiative:         token.Initiative + hc.Initiative,
		hitsAllTargets:     hc.HitsAllTargets,
		capitalShipMissile: hc.CapitalShipMissile,
	}

	if hc.Category == TechCategoryBeamWeapon {
		weaponSlot.weaponType = battleWeaponTypeBeam
		weaponSlot.power = int(float64(weaponSlot.power) * (1 + beamBonus))
	} else if hc.Category == TechCategoryTorpedo {
		weaponSlot.weaponType = battleWeaponTypeTorpedo
	}

	return weaponSlot
}

// return true if this battle will have targets
func (b *battle) hasTargets() bool {
	for _, token := range b.tokens {
		if !token.design.Spec.HasWeapons {
			continue
		}
		// starbases may target a scout, but if they can't move towards it, a battle won't be triggered
		if token.Movement > 0 && b.getTarget(token, b.tokens) != nil {
			return true
		}
	}
	return false
}

// runBattle runs a battle!
func (b *battle) runBattle() *BattleRecord {
	if b.planet != nil {
		log.Info().Msgf("Running a battle at %s involving %d players and %d tokens.", b.planet.Name, len(b.players), len(b.tokens))
	} else {
		log.Info().Msgf("Running a battle at (%.2f, %.2f) involving %d players and %d tokens.", b.position.X, b.position.Y, len(b.players), len(b.tokens))
	}

	moveOrder := b.buildMovementOrder(b.tokens)
	for b.round = 1; b.round <= b.rules.NumBattleRounds; b.round++ {

		// each round we build the SortedWeaponSlots list
		// anew to account for ships that were destroyed
		weaponSlots := b.getSortedWeaponSlots(b.tokens)

		// find new targets
		if b.findMoveTargets(b.tokens) {
			// if we still have targets, process the round
			b.record.recordNewRound()

			// movement is a repeating pattern of 4 movement blocks
			// which we figured out in BuildMovement
			roundBlock := b.round % 4
			for _, token := range moveOrder[roundBlock] {
				b.moveToken(token)
			}

			// iterate over each weapon and fire if they have a target
			for _, weaponSlot := range weaponSlots {
				// find all available targets for this weapon
				targets := b.findTargets(weaponSlot, b.tokens)
				b.fireWeaponSlot(weaponSlot, targets)
			}

			// at the end of this round, regenerate shields
			for _, token := range b.tokens {
				if token.stackShields > 0 && token.isStillInBattle() {
					b.regenerateShields(token)
				}
			}
		} else {
			// no one has targets, we are done
			break
		}
	}

	// record destroyed tokens
	for _, token := range b.tokens {
		if token.quantityDestroyed > 0 {
			b.record.recordDestroyedToken(token, token.quantityDestroyed)
		}
	}

	return b.record
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

// Find all the targets for a weapon
func (b *battle) findTargets(weapon *battleWeaponSlot, tokens []*battleToken) (targets []*battleToken) {
	targets = []*battleToken{}
	attacker := weapon.token
	primaryTarget := attacker.PrimaryTarget
	secondaryTarget := attacker.SecondaryTarget

	var primaryTargets []*battleToken
	var secondaryTargets []*battleToken

	// Find all enemy tokens
	for _, token := range tokens {
		if !token.isStillInBattle() {
			continue
		}
		if !b.willAttack(attacker, b.players[attacker.PlayerNum], token.PlayerNum) {
			continue
		}

		// if we will target this
		if b.willTarget(primaryTarget, token) && weapon.isInRange(token) && weapon.willDamage(token) {
			primaryTargets = append(primaryTargets, token)
		} else if b.willTarget(secondaryTarget, token) && weapon.isInRange(token) && weapon.willDamage(token) {
			secondaryTargets = append(secondaryTargets, token)
		}
	}

	// our list of available targets is all primary and all secondary targets in range
	sort.Slice(primaryTargets, func(i, j int) bool {
		return weapon.getAttractiveness(primaryTargets[i]) < weapon.getAttractiveness(primaryTargets[j])
	})
	sort.Slice(secondaryTargets, func(i, j int) bool {
		return weapon.getAttractiveness(secondaryTargets[i]) < weapon.getAttractiveness(secondaryTargets[j])
	})

	targets = append(targets, primaryTargets...)
	targets = append(targets, secondaryTargets...)
	return targets
}

// findMoveTargets allocates targets for each token in a battle.
func (b *battle) findMoveTargets(tokens []*battleToken) bool {
	hasTargets := false
	for _, token := range tokens {
		token.targetedBy = nil
	}
	for _, token := range tokens {
		if token.Movement == 0 || !token.hasWeapons() || !token.isStillInBattle() {
			continue
		}
		token.moveTarget = b.getTarget(token, tokens)
		if token.moveTarget != nil {
			token.moveTarget.targetedBy = append(token.moveTarget.targetedBy, token)
		}
		hasTargets = token.moveTarget != nil || hasTargets
	}

	return hasTargets
}

// buildMovementOrder builds a list of Movers in this battle.
// Each ship moves in order of mass with heavier ships moving first.
// Ships that can move 3 times in a round move first, then ships that move 2 times, then 1.
func (b *battle) buildMovementOrder(tokens []*battleToken) (moveOrder [4][]*battleToken) {
	// our tokens are moved by mass
	tokensByMass := make([]*battleToken, 0)
	for _, token := range tokens {
		if token.Movement > 0 { // starbases don't move
			tokensByMass = append(tokensByMass, token)
		}
	}
	sort.Slice(tokensByMass, func(i, j int) bool {
		return tokensByMass[i].Mass > tokensByMass[j].Mass
	})

	// each token can move up to 3 times in a round
	// ships that can move 3 times go first, so we loop through the moveNum backwards
	// so that our Movers list has ships that move 3 times first
	for moveNum := 2; moveNum >= 0; moveNum-- {
		// for each block of 4 rounds, add each ship to the movement list if it's supposed to move that round
		for roundBlock := 0; roundBlock < 4; roundBlock++ {
			// add each battle token to the movement for this roundBlock
			for _, token := range tokensByMass {
				// movement is between 2 and 10, so we offset it to fit in our MovementByRound table
				movement := token.Movement

				// see if this token can move on this moveNum (i.e. move 1, 2, or 3)
				if movementByRound[movement-2][roundBlock] > moveNum {
					moveOrder[roundBlock] = append(moveOrder[roundBlock], token)
				}
			}
		}
	}
	return moveOrder
}

// Fire the weapon slot towards its target
func (b *battle) fireWeaponSlot(weapon *battleWeaponSlot, targets []*battleToken) {
	if len(targets) == 0 || !weapon.token.isStillInBattle() {
		// no targets, nothing to do
		return
	}

	switch weapon.weaponType {
	case battleWeaponTypeBeam:
		b.fireBeamWeapon(weapon, targets)
	case battleWeaponTypeTorpedo:
		b.fireTorpedo(weapon, targets)
	}
}

// fire a beam weapon slot at a slice of targets
func (b *battle) fireBeamWeapon(weapon *battleWeaponSlot, targets []*battleToken) {
	attackerShipToken := weapon.token
	damage := weapon.power * weapon.quantity * attackerShipToken.Quantity
	remainingDamage := float64(damage)

	log.Debug().Msgf("%v is attempting to fire at %v targets for a total of %v damage", weapon.token, len(targets), damage)

	for _, target := range targets {
		if !target.isStillInBattle() {
			continue
		}
		if remainingDamage == 0 {
			break
		}
		// skip targets that are out of shields for sappers
		if weapon.damagesShieldsOnly && target.stackShields <= 0 {
			continue
		}

		shields := target.stackShields
		armor := target.armor

		// beam weapon damage reduces by up to 10% over range. So a range 2 weapon is reduced 0% at 0 range, 5% at 1 range, and 10% at 2 range
		distance := weapon.token.getDistanceAway(target.Position)
		if weapon.weaponRange > 0 {
			remainingDamage = remainingDamage * (1 - b.rules.BeamRangeDropoff*float64(distance)/float64(weapon.weaponRange))
		}

		// account for beam defense
		remainingDamage = remainingDamage * (1 - target.beamDefense)

		log.Debug().Msgf("%v fired %v %v(s) at %v (shields: %v, armor: %v, distance: %v, %v@%v damage) for %v (range adjusted to %v)", weapon.token, weapon.quantity, weapon.slot.HullComponent, target, shields, armor, distance, target.Quantity, target.Damage, damage, remainingDamage)

		if remainingDamage > float64(shields) && !weapon.damagesShieldsOnly {
			remainingDamage = remainingDamage - float64(shields)
			target.stackShields = 0

			existingDamage := target.Damage * float64(target.QuantityDamaged)
			remainingDamage += existingDamage

			numDestroyed := int(remainingDamage / float64(armor))
			if numDestroyed >= target.Quantity {
				numDestroyed = target.Quantity
				target.Quantity = 0
				target.quantityDestroyed += numDestroyed
				b.board[target.Position.Y][target.Position.X] -= numDestroyed
				remainingDamage -= float64(armor * numDestroyed)

				target.destroyed = true
				log.Debug().Msgf("%v %v %v(s) hit %v, did %v shield damage and %v armor damage and completely destroyed %v", weapon.token, weapon.quantity, weapon.slot.HullComponent, target, shields, int(remainingDamage)-shields, target)

				b.record.recordBeamFire(b.round, weapon.token, weapon.token.Position, target.Position, weapon.slot.HullSlotIndex, *target, shields, int(remainingDamage)-shields, numDestroyed)
			} else {
				if numDestroyed > 0 {
					target.Quantity -= numDestroyed
					target.quantityDestroyed += numDestroyed
					b.board[target.Position.Y][target.Position.X] -= numDestroyed
				}

				remainingDamage -= float64(armor * numDestroyed)

				if remainingDamage > 0 {
					target.Damage = remainingDamage / float64(target.Quantity)
					target.QuantityDamaged = target.Quantity
					log.Debug().Msgf("%v destroyed %v ships, leaving %v damaged %v@%v damage", weapon.token, numDestroyed, target, target.Quantity, target.Damage)
				}

				b.record.recordBeamFire(b.round, weapon.token, weapon.token.Position, target.Position, weapon.slot.HullSlotIndex, *target, shields, int(remainingDamage)-shields, numDestroyed)
				remainingDamage = 0
			}

		} else {
			target.stackShields -= int(remainingDamage)
			target.stackShields = MaxInt(0, target.stackShields) // make sure we don't go below 0
			log.Debug().Msgf("%v firing %v %v(s) did %v damage to %v shields, leaving %v shields still operational.", weapon.token, weapon.quantity, weapon.slot.HullComponent, remainingDamage, target, target.stackShields)
			b.record.recordBeamFire(b.round, weapon.token, weapon.token.Position, target.Position, weapon.slot.HullSlotIndex, *target, int(remainingDamage), 0, 0)
			// no more damage left, all absorbed by shields
			remainingDamage = 0
		}

		// reset damage for the next target
		if weapon.hitsAllTargets {
			remainingDamage = float64(damage)
		}

		log.Debug().Msgf("%v %v %v(s) has %v remaining dp to burn through %v additional targets.", weapon.token, weapon.quantity, weapon.slot.HullComponent, remainingDamage, len(targets)-1)

		target.damaged = true
	}
}

func (b *battle) getDamageDone(token *battleToken, distance int) int {
	damageDone := 0
	for _, weapon := range token.weaponSlots {
		if !weapon.isInRangeValue(distance) {
			// no damage, skip this weapon
			continue
		}
		switch weapon.weaponType {
		case battleWeaponTypeBeam:
			// TODO: consider sappers + shields
			damage := weapon.power * weapon.quantity * token.Quantity
			if weapon.weaponRange > 0 {
				damage = int(float64(damage) * (1 - b.rules.BeamRangeDropoff*float64(distance)/float64(weapon.weaponRange)))
			}
			damageDone += damage
		case battleWeaponTypeTorpedo:
			// assume every torpedo hits
			damageDone += weapon.power * weapon.quantity * token.Quantity
		}
	}

	return damageDone
}

// Fire a torpedo slot from a ship. Torpedos are different than beam weapons
// A ship will fire each torpedo at its target until the target is destroyed, then
// fire remaining torpedos at the next target.
// Each torpedo has an accuracy rating. That determines if it hits. A torpedo that
// misses still explodes and does 1/8th damage to shields
func (b *battle) fireTorpedo(weapon *battleWeaponSlot, targets []*battleToken) {
	attacker := weapon.token
	// damage is power * number of weapons * number of attackers.
	damage := weapon.power
	accuracy := weapon.accuracy
	numTorpedos := weapon.quantity * attacker.Quantity

	log.Debug().Msgf("%s is attempting to fire at %d targets with %d torpedos at %.2f%% accuracy for %d damage each",
		weapon.token, len(targets), numTorpedos, accuracy*100.0, damage)

	// fire each torpedo at each target until it's destroyed or we're out of torpedos
	remainingTorpedos := numTorpedos
	torpedoNum := 0
	for _, target := range targets {
		if !target.isStillInBattle() {
			// this token isn't valid anymore, skip it
			continue
		}

		// no more damage to spread, break out
		if remainingTorpedos == 0 {
			break
		}

		// shields are shared among all tokens
		armor := target.armor
		shipDamage := target.Damage

		totalShieldDamage := 0
		totalArmorDamage := 0
		hits := 0
		misses := 0
		shipsDestroyed := 0

		for remainingTorpedos > 0 && !target.destroyed {
			// fire a torpedo
			torpedoNum++
			remainingTorpedos--
			hit := (accuracy - target.torpedoJamming) > float64(b.rules.random.Float64())

			if hit {
				hits++

				// torpedos do half damage to shields, half to armor (until shields are gone, when they do full armor damage)
				shieldDamage := float64(0.5) * float64(damage)
				armorDamage := float64(0.5) * float64(damage)

				// apply up to half our damage to shields
				// anything leftover goes to armor
				afterShieldsDamaged := float64(target.stackShields) - shieldDamage
				var actualShieldDamage float64
				if afterShieldsDamaged < 0 {
					// We did more damage to shields than they had remaining
					// apply the difference to armor
					actualShieldDamage = shieldDamage + afterShieldsDamaged
					armorDamage += float64(-afterShieldsDamaged)

				} else {
					actualShieldDamage = shieldDamage
				}
				target.stackShields -= int(actualShieldDamage)

				if target.stackShields <= 0 && weapon.capitalShipMissile {
					// capital ship missiles double damage after shields are gone
					armorDamage *= 2
				}

				totalShieldDamage += int(actualShieldDamage)
				totalArmorDamage += int(armorDamage)
				shipDamage += armorDamage

				// this torpedo blew up a ship, hooray!
				if shipDamage >= float64(armor) {
					// remove a ship from this stack
					target.Quantity--
					target.quantityDestroyed++
					b.board[target.Position.Y][target.Position.X] -= 1
					target.QuantityDamaged = MaxInt(target.QuantityDamaged-1, 0)

					if target.QuantityDamaged > 0 {
						// we destroyed a token, but we still have damaged tokens in the stack
						// so reset our shipDamage counter to the damage + any leftover. We apply that
						// to the rest of the tokens
						// i.e. if we fire 2 omega torpedos for 300 damage each at 3 damaged 1700dp@1300 ships
						// the first shot damages the top ship, the second one kills it but we have 200 leftover
						// this will carry over to damage the remaining ships
						leftoverDamage := shipDamage - float64(armor)
						shipDamage = target.Damage + leftoverDamage
					} else {
						// we have no more damaged tokens, so remove the stack's damage
						// and reset our ship damage to 0
						// this could happen if we are firing on a stack with 3 ships but 2@10 damage or something
						shipDamage = 0
						target.Damage = 0
					}
					if target.Quantity <= 0 {
						// record that we destroyed this token
						target.destroyed = true
						log.Debug().Msgf("%v torpedo number %v hit %v, did %v shield damage and %v armor damage and completely destroyed %v", weapon.token, torpedoNum, target, actualShieldDamage, armorDamage, target)
						shipsDestroyed++
					}
				} else {
					log.Debug().Msgf("%v torpedo number %v hit %v, did %v shield damage and %v armor damage (%v accumulated damage so far)", weapon.token, torpedoNum, target, actualShieldDamage, armorDamage, shipDamage)
				}
			} else {
				misses++
				// damage shields by 1/8th
				// round up, do a minimum of 1 damage
				shieldDamage := int(math.Min(1, math.Round(b.rules.TorpedoSplashDamage*float64(damage))))
				actualShieldDamage := shieldDamage
				if shieldDamage > target.stackShields {
					actualShieldDamage = target.stackShields
				}
				target.stackShields = int(math.Max(0, float64(target.stackShields-shieldDamage)))
				log.Debug().Msgf("%s torpedo number %d missed %s, did %d damage to shields leaving %d shields", weapon.token, torpedoNum, target, shieldDamage, target.stackShields)

				totalShieldDamage += actualShieldDamage
			}
		}

		// we have leftover damage, apply it to all remaining tokens evenly
		if shipDamage > 0 && target.Quantity > 0 {
			target.damaged = true // target lived, but is damaged
			var previousDamage float64
			if target.QuantityDamaged > 0 {
				// we had some tokens damaged previously that we didn't touch
				previousDamage = target.Damage * float64(target.QuantityDamaged)
				shipDamage -= target.Damage // we already include this in our ship damage
			}
			target.Damage = (shipDamage + previousDamage) / float64(target.Quantity)
			target.QuantityDamaged = target.Quantity
			log.Debug().Msgf("%s had %d hits and %d misses to %v for %v total damage leaving %d@%v", weapon.token, hits, misses, target, totalArmorDamage+totalShieldDamage, target.QuantityDamaged, target.Damage)

		}
		b.record.recordTorpedoFire(b.round, weapon.token, weapon.token.Position, target.Position, weapon.slot.HullSlotIndex, target, totalShieldDamage, totalArmorDamage, shipsDestroyed, hits, misses)
	}
}

// moveToken moves a token towards or away from its target
func (b *battle) moveToken(token *battleToken) {
	if !token.isStillInBattle() {
		return
	}
	// count this token's moves
	token.movesMade++

	// always disengage if we have no weapons or targets
	if token.moveTarget == nil || !token.hasWeapons() {
		token.Tactic = BattleTacticDisengage
	}

	// disengage if we are a scared little token
	if token.Tactic == BattleTacticDisengageIfChallenged && token.damaged {
		token.Tactic = BattleTacticDisengage
	}

	oldPosition := token.Position

	if token.Tactic == BattleTacticDisengage {
		b.runAway(token)
	} else {
		// create a move record for the viewer and then move the token
		bestMove := b.getBestMove(token)
		b.record.recordMove(b.round, token, token.Position, bestMove)
		token.Position = bestMove
	}

	// update the board after a token moves
	b.board[oldPosition.Y][oldPosition.X] -= token.Quantity
	b.board[token.Position.Y][token.Position.X] += token.Quantity
}

// regenerateShields regenerates the shields of the given token if the player regenerates shields
// and the token has shields.
func (b *battle) regenerateShields(token *battleToken) {
	player := b.players[token.PlayerNum]

	if player.Race.Spec.ShieldRegenerationRate > 0 && token.stackShields > 0 {
		regenerationAmount := int(float64(token.totalStackShields)*player.Race.Spec.ShieldRegenerationRate + 0.5)
		token.stackShields = int(Clamp(token.stackShields+regenerationAmount, 0, token.totalStackShields))
	}
}

// for a new position on the board, and an existing slice of bestMoves, return the new bestMoves
// if better, take the new position and discard the others
// if not better, but closer to center, take the new move and discard the others
// if not better and farther away, keep the bestMoves as is
// if equivalent and the same distance from center, add it to bestMoves
func updateBestMoves(better bool, newPosition BattleVector, bestMoves []BattleVector) []BattleVector {
	// center of battle board is at (4.5, 4.5) so scale to 100x100 for 45,45 scaled center
	newDistance := newPosition.scale(10).distance(BattleVector{45, 45})
	oldDistance := bestMoves[0].scale(10).distance(BattleVector{45, 45})

	if better || newDistance < oldDistance {
		// this move is strictly better than other moves based on damage or proximity to center
		bestMoves = nil
	} else if oldDistance < newDistance {
		// damage is the same but proximity is worse
		return bestMoves
	}
	return append(bestMoves, newPosition)
}

// get the best move to attack a target
// this checks all the nearby squares and moves in the direction that gets the token
// closer or does the best damage, according to the battle plan
// i.e. max damage, max damage ratio, max net damage, etc
func (b *battle) getBestMove(token *battleToken) BattleVector {

	bestDamageDone := 0
	bestDamageTaken := math.MaxInt
	bestDamageRatio := 0.0
	bestNetDamage := -math.MaxInt
	bestMoves := make([]BattleVector, 0, 9)
	bestMoves = append(bestMoves, token.Position)
	lowestDamageMoves := make([]BattleVector, 0, 9)
	lowestDamageMoves = append(lowestDamageMoves, token.Position)
	bestDamageTakenForLowestDamageMove := math.MaxInt

	// target is either whoever this token is targeting, or targeted by
	target := token.moveTarget
	var currentDistance int
	if target != nil {
		currentDistance = token.getDistanceAway(target.Position)
	} else {
		currentDistance = math.MaxInt
		for _, targetedBy := range token.targetedBy {
			// find the closest ship targeting us
			// TODO: maybe find the scariest target?
			dist := token.getDistanceAway(targetedBy.Position)
			if dist < currentDistance {
				currentDistance = dist
				target = targetedBy
			}
		}
	}

	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			newPosition := BattleVector{token.Position.X + dx, token.Position.Y + dy}
			if newPosition.X < 0 || newPosition.X >= battleWidth || newPosition.Y < 0 || newPosition.Y >= battleHeight {
				// skip invalid squares
				continue
			}

			// if we are running away and no one is targeting us
			if target == nil {
				// no one is targeting us, add this square to random moves
				lowestDamageMoves = append(lowestDamageMoves, newPosition)
				continue
			}

			// figure out how well this move is
			distance := newPosition.distance(target.Position)
			damageDone := b.getDamageDone(token, distance)
			damageTaken := 0
			for _, attacker := range token.targetedBy {
				distanceToAttacker := newPosition.distance(attacker.Position)
				damageTaken += b.getDamageDone(attacker, distanceToAttacker)
			}

			// if this will move us closer to our target, see if it's the best low damage move
			// or, if this distance is greater than our currentDistance and we're running away, see if it's the best low damage move
			if distance < currentDistance && bestDamageTakenForLowestDamageMove >= damageTaken ||
				token.Tactic == BattleTacticDisengage && distance > currentDistance && bestDamageTakenForLowestDamageMove >= damageTaken {
				lowestDamageMoves = updateBestMoves(bestDamageTakenForLowestDamageMove > damageTaken, newPosition, lowestDamageMoves)
				bestDamageTakenForLowestDamageMove = damageTaken
			}

			switch token.Tactic {
			case BattleTacticMaximizeDamageRatio:
				if damageTaken == 0 {
					if damageDone >= bestDamageDone {
						bestMoves = updateBestMoves(damageDone > bestDamageDone, newPosition, bestMoves)
						bestDamageTaken = 0
						bestDamageDone = damageDone
					}
				} else if bestDamageDone == 0 && float64(damageDone)/float64(damageTaken) >= bestDamageRatio {
					bestMoves = updateBestMoves(float64(damageDone)/float64(damageTaken) > bestDamageRatio, newPosition, bestMoves)
					bestDamageRatio = float64(damageDone) / float64(damageTaken)
				}
			case BattleTacticMaximizeNetDamage:
				if damageDone > 0 && damageDone-damageTaken >= bestNetDamage {
					bestMoves = updateBestMoves(damageDone-damageTaken > bestNetDamage, newPosition, bestMoves)
					bestDamageDone = damageDone
					bestNetDamage = damageDone - damageTaken
				}

			case BattleTacticMinimizeDamageToSelf:
				if damageTaken <= bestDamageTaken {
					bestMoves = updateBestMoves(damageTaken < bestDamageTaken, newPosition, bestMoves)
					bestDamageTaken = damageTaken
				}

			case BattleTacticMaximizeDamage:
				if damageDone >= bestDamageDone {
					bestMoves = updateBestMoves(damageTaken > bestDamageTaken, newPosition, bestMoves)
					bestDamageDone = damageDone
				}
			}
		}
	}

	// if we are getting away, or none of our moves lead to damage, pick the lowest damage move
	if token.Tactic == BattleTacticDisengage || (bestDamageDone == 0 && bestDamageRatio == 0.0) {
		return lowestDamageMoves[b.rules.random.Intn(len(lowestDamageMoves))]
	}

	return bestMoves[b.rules.random.Intn(len(bestMoves))]
}

func (b *battle) runAway(token *battleToken) {

	if token.movesMade >= b.rules.MovesToRunAway {
		// we've moved enough to leave the board
		token.ranAway = true
		b.board[token.Position.Y][token.Position.X] -= token.Quantity
		b.record.recordRunAway(b.round, token)
		return
	}

	// find the best move for running away
	newPosition := b.getBestMove(token)
	b.record.recordMove(b.round, token, token.Position, newPosition)
	token.Position = newPosition
}

// getTarget returns the primary or secondary target based on the attacker's battle plan and the defenders available
func (b *battle) getTarget(attacker *battleToken, defenders []*battleToken) *battleToken {
	var primaryTarget *battleToken
	var secondaryTarget *battleToken
	var primaryTargetAttractiveness float64
	var secondaryTargetAttractiveness float64
	primaryTargetOrder := attacker.PrimaryTarget
	secondaryTargetOrder := attacker.SecondaryTarget

	// TODO: We need to account for the fact that if a fleet targets us, we will target them back
	for _, defender := range defenders {
		if !defender.isStillInBattle() {
			continue
		}
		if b.willAttack(attacker, b.players[attacker.PlayerNum], defender.PlayerNum) {
			// if we would target this defender with our primary target and it's more attractive than our current primaryTarget, pick it
			if b.willTarget(primaryTargetOrder, defender) {
				attractiveness := defender.getAttractiveness(attacker)
				if attractiveness >= primaryTargetAttractiveness {
					primaryTarget = defender
					primaryTargetAttractiveness = attractiveness
				}
			}

			// if we would target this defender with our secondary target, pick it
			if b.willTarget(secondaryTargetOrder, defender) {
				attractiveness := defender.getAttractiveness(attacker)
				if attractiveness >= secondaryTargetAttractiveness {
					secondaryTarget = defender
					secondaryTargetAttractiveness = attractiveness
				}
			}
		}
	}

	if primaryTarget != nil {
		return primaryTarget
	}

	return secondaryTarget
}

// get the attractiveness of a token based on cost and defense capabilities
func (defender *battleToken) getAttractiveness(attacker *battleToken) float64 {
	cost := defender.design.Spec.Cost.MultiplyInt(defender.Quantity)
	defense := (defender.armor + defender.shields) * defender.Quantity

	return float64(float64(cost.Germanium+cost.Resources) / float64(defense))
}

// get the attractiveness of a token versus a weapon
func (weapon *battleWeaponSlot) getAttractiveness(target *battleToken) float64 {
	cost := target.design.Spec.Cost.MultiplyInt(target.Quantity)
	defense := float64((target.armor + target.shields) * target.Quantity)

	// increase the defense for jammers and beam deflectors
	switch weapon.weaponType {
	case battleWeaponTypeBeam:
		defense = defense * (1 + target.beamDefense)
	case battleWeaponTypeTorpedo:
		defense = defense * (1 + target.torpedoJamming)
	}

	return float64(float64(cost.Germanium+cost.Resources) / float64(defense))
}

// return true if this fleet will attack a fleet by another player based on player
// relations and the fleet battle plan
func (b *battle) willAttack(token *battleToken, player *Player, otherPlayerNum int) bool {
	// if we have weapons and we don't own this other fleet, see if we
	// would target it
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

// willTarget returns true if the BattleOrder Target type would target this token
func (b *battle) willTarget(target BattleTarget, token *battleToken) bool {
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

// get all weapon slots on the board, sorted by initiative
func (b *battle) getSortedWeaponSlots(tokens []*battleToken) []*battleWeaponSlot {
	slots := []*battleWeaponSlot{}
	for _, token := range tokens {
		if token.isStillInBattle() {
			slots = append(slots, token.weaponSlots...)
		}
	}
	sort.Slice(slots, func(i, j int) bool {
		return slots[i].initiative > slots[j].initiative
	})
	return slots
}

// Function to allow server to run test battles
func RunTestBattle(players []*Player, fleets []*Fleet) *BattleRecord {
	rules := NewRules()

	playersByNum := map[int]*Player{}
	designsByNum := make(map[playerObject]*ShipDesign)
	battlePlansByName := make(map[playerBattlePlanNum]*BattlePlan)

	for _, player := range players {
		playersByNum[player.Num] = player
		player.Race.Spec = computeRaceSpec(&player.Race, &rules)
		player.Spec = computePlayerSpec(player, &rules, []*Planet{})

		for _, design := range player.Designs {
			design.Spec = ComputeShipDesignSpec(&rules, player.TechLevels, player.Race.Spec, design)
			designsByNum[playerObjectKey(design.PlayerNum, design.Num)] = design
		}

		for i := range player.BattlePlans {
			plan := &player.BattlePlans[i]
			battlePlansByName[playerBattlePlanNum{PlayerNum: player.Num, Num: plan.Num}] = plan
		}

	}

	for _, fleet := range fleets {
		for i := range fleet.Tokens {
			token := &fleet.Tokens[i]
			token.design = designsByNum[playerObjectKey(fleet.PlayerNum, token.DesignNum)]
		}
		fleet.Spec = ComputeFleetSpec(&rules, playersByNum[fleet.PlayerNum], fleet)
		fleet.battlePlan = battlePlansByName[playerBattlePlanNum{fleet.PlayerNum, fleet.BattlePlanNum}]
	}

	battler := newBattler(&rules, &StaticTechStore, 1, playersByNum, fleets, nil)
	record := battler.runBattle()
	for _, player := range players {

		for _, otherplayer := range players {
			player.discoverer.discoverPlayer(otherplayer)
		}
		for _, fleet := range fleets {
			if fleet.PlayerNum != player.Num {
				for _, token := range fleet.Tokens {
					player.discoverer.discoverDesign(token.design, true)
				}
			}
		}
	}

	return record
}
